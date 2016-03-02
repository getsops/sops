# Load in our dependencies
import json

# Define our configuration
# DEV: THE FOLLOWING CONFIGURATIONS SHOULD NOT CONTAIN ANY SECRETS
#   THIS FILE IS NOT ENCRYPTED!!
common = {
    'port': 8080,
}

development = {

}

test = {

}

production = {

}

config = {
    'common': common,
    'development': development,
    'test': test,
    'production': production,
}


def walk(item, fn):
    """Traverse dicts and lists to update keys via `fn`"""
    # If we are looking at a dict, then traverse each of its branches
    if isinstance(item, dict):
        for key in item:
            # Walk our value
            walk(item[key], fn)

            # If we are changing our key, then update it
            new_key = fn(key)
            if new_key != key:
                item[new_key] = item[key]
                del item[key]
    # Otherwise, if we are looking at a list, walk each of its items
    elif isinstance(item, list):
        for val in item:
            walk(val, fn)


# Merge all of our static secrets onto our config
# For each of our secrets
secret_files = [
    'config/static_github.json',
]
for secret_file in secret_files:
    with open(secret_file, 'r') as file:
        # Load and parse our JSON
        data = json.loads(file.read())

        # Strip off `_unencrypted` from all keys
        walk(data, lambda key: key.replace('_unencrypted', ''))

        # For each of the environments
        for env_key in data:
            # Load in the respective source and target
            env_src = data[env_key]
            env_target = config[env_key]

            # Merge info between configs
            for key in env_src:
                if key in env_target:
                    raise AssertionError(
                        'Expected "{env_key}.{key}" to not be defined already '
                        'but it was. '
                        'Please verify no configs are using the same key'
                        .format(env_key=env_key, key=key))
            env_target.update(env_src)
