# Bitbucket Cloud Pipeline Runner

This is my first attempt a developing a CLI tool using Golang, so please bare with me.

The runner has been developed because of limitation within BitBucket pipelines, there is no native support in the YAML
spec for triggering other piplines and the
[Trigger pipeline](https://bitbucket.org/atlassian/trigger-pipeline/src/master/) pipe whilsts fills the lack of the
native support, it can add an incredible about of noise to your pipelines and doesn't send the output to the triggering
pipeline, you must click through to it.

The first iteration is intended to provide something like the previously mentioned pipe so the implementation can be
proven. Then support for having pipeline configuration stored in files.

## Notes

- Please try not to judge the code too much. I'm new to golang and I just want to get it working first.
- Tests are lacking at the moment

## Configuration

Regardless of the configuration choice, you must have an App Password setup with `write` access to
`pipelines`.

### ENV Based

```bash
export BPR_BITBUCKET_USERNAME="username"
export BPR_BITBUCKET_PASSWORD="password"
```

### File based

Bpr will automatically look in your home directory for `./bpr/config.env`. However, you can override this using
the `BPR_CONFIG_PATH` environment variable which must be an absolute path to the config.

```ini
; ~/.bpr/config.env
BITBUCKET_USERNAME=username
BITBUCKET_PASSWORD=password
```

## Usage

| FLAG     | DESCRIPTION                                |
| -------- | ------------------------------------------ |
| owner    | the org/workspace/user that owns the repo  |
| ref      | git branch the custom pipeline lives under |
| pipeline | the name of the custom pipe                |
| vars     | A list of variables provided in JSON       |

`owner`, `ref`, and `pipeline` must all be used together.

### Example

#### Example 1 - Flags Only - Currenctly implemented

```bash
bpr -owner 'DanielChalk' \
    -repo 'bitbucket-pipeline-runner' \
    -ref 'wip' \
    -pipeline 'example' \
    -vars '[{"key":"NAME", "value":"daniel"}]'
```

#### Example 2 - File based

```yaml
# pipeline-spec.yml
pipelines:
  my_pipeline:
    pipeline: workspace/repo/branch/pipeline
    variables: 
      KEY: Value
```

```bash
bpr pipeline-spec.yml
```

#### Example 3 - File + Vars - to be implemented

This method allows you to have a generic file with most of the configuration, but the `-vars` flag allows you to
override and append values.

```yaml
# pipeline-spec.yml
pipelines:
  my_pipeline:
    pipeline: workspace/repo/branch/pipeline
    # No vars we we will set them externally
```

```bash
bpr \
  -vars '[{"key": "NAME", "value": "Daniel"}]' \
  pipeline-spec.yml
```

### Variables

Variables are set in two different schema types based on whether they are set via flags, or a YAML file.

#### YAML

In the YAML spec files, variables are simple key/value pairs. This is because we do not want to allow secured variables
being committed in them.

```yaml
# You can also set global variables for the pipelines in your specs file, making for less copying and pasting.
variables:
  GLOBAL_VAR: var
pipelines: 
  my_pipeline:
    pipeline: workspace/repo/branch/pipeline
    # Variables are key/value only, if you want to use secured variables, you must provide them as a flag    
    variables:
      KEY: Value
```

#### Flag

The `-vars` command line flag takes a JSON list of variables, which support the "secured" property

```bash
$VARIABLES='[{"key": "MY_VAR_NAME", "value": "MY_VAR_VALUE" }, {"key": "MY_SECURE_VAR_NAME", "value": "MY_SECURE_VAR_VALUE", "secured": true}]'
bpr -vars $VARIABLES
```

#### Variable precedence

Variables are layered onto each other in the following order:

- Global spec file variables (the lowest precedence)
- Spec pipeline variables
- Variables provided by the flag (the highest precedence)