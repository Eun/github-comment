# github-comment

Create or update an github comment

Usage: 
```bash
# Create a new comment in the pull request owner/repo/pulls/2 wth the unique id 123-ABC
github-comment --repo owner/repo --pr 2 --id "123-ABC" "Hello World"

# Update the comment
github-comment --repo owner/repo --pr 2 --id "123-ABC" "Hello Universe"

```


