# Bitbucket Cloud Pipeline Runner

This is my first attempt a developing a CLI tool using Golang, so please bare with me.

The runner has been developed because of limitation within BitBucket pipelines, there is no native support in the YAML spec for triggering other piplines and the 
[Trigger pipeline](https://bitbucket.org/atlassian/trigger-pipeline/src/master/) pipe whilsts fills the lack of the native support, it can add an incredible about of noise to your pipelines and doesn't send the output to the triggering pipeline, you must click through to it.

The first iteration is intended to provide something like the previously mentioned pipe so the implementation can be proven. Then support for having pipeline configuration stored in files.

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

Bpr will automatically look in your home directory for `./bpr/config.env`. However, you can override
this using the `BPR_CONFIG_PATH` environment variable which must be an absolute path to the config.


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
- workspace: DanielChalk
  repo_slug: bitbucket-pipeline-runner
  ref: wip
  pipeline: example
  variables: 
    - key: NAME
      value: Daniel Chalk
```

```bash
bpr pipeline-spec.yml
```

#### Example 3 - File + Vars - to be implemented

This method allows you to have a generic file with most of the configuration, but the `-vars` flag
allows you to override and append values.

```yaml
# pipeline-spec.yml
- workspace: DanielChalk
  repo_slug: bitbucket-pipeline-runner
  ref: wip
  pipeline: example
  # No vars we we will set them externally
```

```bash
bpr \
  -vars '[{"key": "NAME", "value": "Daniel"}]' \
  pipeline-spec.yml
```

### Variables

You must format variables according to the following specification

```json
[
    { 
        "key": "MY_VAR_NAME", 
        "value": "MY_VAR_VALUE",
    },
    { 
        "key": "MY_SECURE_VAR_NAME", 
        "value": "MY_SECURE_VAR_VALUE",
        "secure": true
    }
]
```

As you can see you can provide secure params which are masked in the output, but the will exist, in
your command history or even the pipelines you hardcode its usage in.