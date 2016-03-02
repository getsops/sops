# Load in our dependencies
from __future__ import absolute_import
from config.static import config


# Define our main function
def main():
    # Output our configuration
    # DEV: We use custom keys for a custom sort
    print('Configuration')
    print('=============')
    for env_key in ('common', 'development', 'test', 'production'):
        # Example: `Environment: common`
        # Example: `-------------------`
        env_str = 'Environment: {env_key}'.format(env_key=env_key)
        print(env_str)
        print(''.join(['-' for char in env_str]))

        env_config = config[env_key]
        for key in sorted(env_config.keys()):
            # Example: `port: "8080"`
            print('{key}: "{val}"'.format(key=key, val=env_config[key]))
        print('')


# If this script is being invoked directly, then run our main function
if __name__ == '__main__':
    main()
