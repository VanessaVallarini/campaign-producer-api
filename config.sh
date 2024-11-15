#!/bin/bash

# Check if environment was provided
if [ -z "$1" ]; then
  echo "Error: No environment provided. Please provide the environment name (ex: canary, dev, prod)."
  exit 1
fi

environment=$1
overrideEnv=$environment

if [[ "$environment" == "local" ]]; then
  environment="sandbox"
fi

# The path to the yaml file
yamlFile="./k8s/values-$environment.yaml"
configPath="./internal/config"
localPath="./local"
secretYaml="$localPath/secret-$overrideEnv.yml"
envOverrideYaml="$localPath/override-$overrideEnv.yml"
fileContent=""

# Check if both secretYaml and yq are available
if [[ ! -f "$secretYaml" ]] || [[ ! -s "$secretYaml" ]]; then
  read -n1 -r -p $'File secret not found! you are sure that want continue?\nPress any key to continue...\n or ctrl + c to stop' key
fi

# Check if both yamlFile and yq are available
if [[ -f "$yamlFile" ]] && [[ -x "$(command -v yq)" ]]; then
  if [[ "$environment" == "production" ]]; then
    fileContent=$(yq e '.appConfigFile.data[0].fileContent' "$yamlFile")
  else
    fileContent=$(yq e '.appConfigFile.data[0].fileContent' "$yamlFile")
  fi
  fileContent=$(yq eval-all 'select(fileIndex == 0) * select(fileIndex == 1)' <(echo "$fileContent") <(cat "$secretYaml"))
  fileContent=$(yq eval-all 'select(fileIndex == 0) * select(fileIndex == 1)' <(echo "$fileContent") <(cat "$envOverrideYaml"))
else
  echo "File $yamlFile does not exist or yq is not installed."
fi

echo -e "$fileContent" > "$configPath/configuration.yml"
