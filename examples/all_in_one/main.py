# Load in our dependencies
from __future__ import absolute_import
from config.static import common


# Define our main function
def main():
    # Output our configuration
    print('Configuration')
    print('-------------')
    for key in common:
        # Example: `port: "8080"`
        print('{key}: "{val}"'.format(key=key, val=common[key]))


# If this script is being invoked directly, then run our main function
if __name__ == '__main__':
    main()
