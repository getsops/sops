import unittest2
import mock
import sys

import sops

if sys.version_info[0] == 2:
    import __builtin__ as builtins  # pylint:disable=import-error
else:
    import builtins  # pylint:disable=import-error


class TreeTest(unittest2.TestCase):

    def test_json_loader_is_used_on_json_filetype(self):
        # XXX put some real json here.
        m = mock.mock_open(read_data='"content"')
        with mock.patch.object(builtins, 'open', m):
            assert sops.load_tree('path', 'json') == "content"

    def test_yaml_loader_is_used_on_yaml_filetype(self):
        # XXX put some real yaml here.
        m = mock.mock_open(read_data='"content"')
        with mock.patch.object(builtins, 'open', m):
            assert sops.load_tree('path', 'yaml') == "content"

    @mock.patch('sops.json.load')
    def test_example_with_a_mocked_call(self, json_mock):
        m = mock.mock_open(read_data='"content"')
        with mock.patch.object(builtins, 'open', m):
            sops.load_tree('path', 'json')
            json_mock.assert_called_with(m())
