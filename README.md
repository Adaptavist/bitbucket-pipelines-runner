# Bitbucket Cloud Pipeline Runner

This is my first attempt a developing a CLI tool using Golang, so please bare with me.

The runner has been developed because of limitation within BitBucket pipelines, there is no native 
support in the YAML spec for triggering other piplines, while
the [Trigger pipeline](https://bitbucket.org/atlassian/trigger-pipeline/src/main/) pipe fills the
lack of the native support, it can add an incredible amount of noise to your pipelines and doesn't
send the output to the triggering pipeline, you must click through to it.

## Notes

- Please try not to judge the code too much. I'm new to golang and I just want to get it working first.
- Tests need work, we've got some but not enough.

## Configuration

Regardless of the configuration choice, you must have an App Password setup with `write` access to
`pipelines` and `administrator` on `repositories`.

**Using env vars**

```bash
export BPR_BITBUCKET_USERNAME="username"
export BPR_BITBUCKET_PASSWORD="password"
```

**Using config file**

```ini
; ~/.bpr.env
BITBUCKET_USERNAME="username"
BITBUCKET_PASSWORD="password"
```

## Commands

### Pipeline

Runs a single pipeling via flags

`bpr pipeline $workspace/$repo_slug/$ref_type/$ref_name[/pipeline_name]`

- `--var 'key=value'` a repeatable flag for providing variables to your pipeline
- `--secrete 'key=value'` a repeatable flag for providing secured variables to your pipeling
- `--dry` shows what the command will do rather than do it

#### Examples

```bash
# runs default pipeline on the main branch
bpr pipeline Owner/repo-slug/branch/main
# runs a custom pipeline (my-pipeline) on main
bpr pipeline Owner/repo-slug/branch/main/my-pipeline
# runs a custom pipeline (my-pipeline) on a tag
bpr pipeline Owner/repo-slug/tag/v1.0.0/my-pipeline
# runs a pipeline with variables/secrets to send to your pipeline
bpr pipeline Owner/repo-slug/branch/main --var 'username=user' --secret 'password=password' --var 'timeout=10'
```

### Spec

Loads all `.bpr.yml` files in your current working directory to build a list of pipelines to run. Then either runs all
of pipelines or just one if you provide the `--only` flag. Secrets are deliberately left our of the YAML spec, as we
don't want to risk of them being introduced into source control.

#### Example Spec file

```yaml
# Variables global to pipelines created with in
variables:
  USERNAME: username
pipelines:
  # Pipeline keys are be globally unique to your working directory, not just the file
  my_pipeline_key:
    pipeline: Owner/repo-slug/branch/example # defaults to "default"
  my_other_pipeliny_key:
    pipeline: Owner/repo-slug/branch/example/pipeline # custom pipeline on main
    variables:
      WAIT: 10 # Provides WAIT as a variable to the pipeline
```

#### Running Specs

Using spec files with bpr requires your current working directory to be where your `.bpr.yml` are.

Run everything:

```bash
bpr spec
```

Run a specific pipeline by its key in the spec:

```bash
bpr spec --only "my_pipeline"
```

Run piplines with additional variables, works the same as when running pipelines directly

```bash
bpr spec --var 'key=value' --secret 'key2=value'
```

Do a dry run

```bash
bpr spec --dry
```

You can also override the target branch, tag, and pipeline. Use this with great care as it will override the target for
all pipelines, so its recomended to only be used with the `--only` flag.

```bash
bpr spec --target-type "tag" --target-ref "main" --target-pipeline "custom_pipeline"
```
