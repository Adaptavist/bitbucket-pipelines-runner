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

`bpr [args] workspace/repo_slug/branch[/pipeline_name]`

- `--var 'key=value'` a repeatable flag for providing variables to your pipeline
- `--secrete 'key=value'` a repeatable flag for providing secured variables to your pipeling
- `--dry` shows what the command will do rather than do it

#### Examples

```bash
# runs default pipeline on the master branch
bpr pipeline Owner/repo-slug/master 
# runs a custom pipeline (my-pipeline) on master
bpr pipeline Owner/repo-slug/master/my-pipeline
# runs a pipeling with variables/secrets to send to your pipeline
bpr pipeline Owner/repo-slug/master --var 'username=user' --secret 'password=password' --var 'timeout=10'
```

### Spec

Loads all `.bpr.yml` files in your current working directory to build a list of pipelines to run.

#### Example Spec file

```yaml
# Variables global to pipelines created with in
variables:
  USERNAME: username
pipelines:
  # Pipeline keys are be globally unique to your working directory, not just the file
  my_pipeline_key:
    pipeline: Owner/repo-slug/branch # defaults to "default" 
  my_other_pipeliny_key:
    pipeline: Owner/repo-slug/branch/pipeline # custom pipeline on master
    variables: 
      WAIT: 10 # Provides WAIT as a variable to the pipeline
```

#### Running Spec FilterSteps

```bash
cd my-pipeline-dirs
# Run all the pipelines found in your directory
bpr spec 
# Run a specific pipeline
bpr spec --only my_pipeline
# Run piplines with additional variables
bpr spec --var 'key=value' --secret 'key2=value'
# Do a dry run
bpr spec --dry
```
