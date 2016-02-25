# Load in our dependencies
import json

# Load in our secrets
with open('config/secret.json', 'r') as file:
    secret = json.loads(file.read())

# Define our configuration
common = {
    'github_oauth_token': secret['github_oauth_token'],
    'port': 8080,
}
