## Jira-Github-Release-Automation-with-go

Github actions for Jira release automation.

### Inputs

`domain`: Domain Name (check your jira url) - https://{domain}.atlassian.net
`project`: Your project id
`releaseName`: name of release you want
`auth-token`: auth-token key
`create-next-version`: create next version (option)

### Auth Token

https://developer.atlassian.com/cloud/jira/platform/basic-auth-for-rest-apis/

1. Generate an API token for Jira using your Atlassian Account.
2. Build a string of the form useremail:api_token. (ted@prnd.co.kr:xxxxxxx)
3. BASE64 encode the string.

- Linux/Unix/MacOS:

```shell
echo -n user@example.com:api_token_string | base64
```

- Windows 7 and later, using Microsoft Powershell:

```shell
$Text = ‘user@example.com:api_token_string’
$Bytes = [System.Text.Encoding]::UTF8.GetBytes($Text)
$EncodedText = [Convert]::ToBase64String($Bytes)
$EncodedText
```

### Example Usage

```yaml
name: Jira Release
on:
  push:
    branches: [main]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Extract version name
        run: echo "##[set-output name=version;]$(echo '${{ github.event.head_commit.message }}' | egrep -o '[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}')"
        id: extract_version_name
      - name: Jira Github Release Automation
        id: release
        uses: Geun-Oh/Jira-Github-Release-Automation-With-Go@0.1.0
        with:
          domain: "My domain"
          project: "My project"
          releaseName: "My release version name"
          auth-token: "My release token"
          create-next-version: false ## can be true
      - name: Print New Version ## when create-next-version option is true
        run: |
          echo ${{ steps.release.outputs.new-version }}
```

### Ref

https://github.com/PRNDcompany/jira-release
