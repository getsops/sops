# -*- coding: utf-8 -*-
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
#
# Contributor: Julien Vehent <jvehent@mozilla.com> [:ulfr]
# Contributor: Daniel Thornton <daniel@relud.com>
# Contributor: Alexis Metaireau <alexis@mozilla.com> [:alexis]
# Contributor: Rémy Hubscher <natim@mozilla.com> [:natim]

import unittest2
import mock
import sys

import sops

if sys.version_info[0] == 2:
    import __builtin__ as builtins
else:
    import builtins


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

    def test_detect_filetype_handle_json(self):
        assert sops.detect_filetype("file.json") == "json"

    def test_detect_filetype_handle_yml(self):
        assert sops.detect_filetype("file.yml") == "yaml"

    def test_detect_filetype_handle_yaml(self):
        assert sops.detect_filetype("file.yaml") == "yaml"

    def test_detect_filetype_returns_text_if_unknown(self):
        assert sops.detect_filetype("file.xml") == "text"

    def test_verify_or_create_sops_branch(self):
        """Verify or create the sops branch"""
        # - sops is created if missing from tree
        # - kms arn is used
        # - pgp fp is used
        # - SOPS_KMS_ARN env variable is used
        # - SOPS_PGP_FP env variable is used
        # - panic error is raise and program quit with code 111 if
        #   nothing is defined

    # Test decryption
    def test_walk_and_decrypt(self):
        """Walk the branch recursively and decrypt leaves."""
        # - test stash value
        # - test dict
        # - test list
        # - test ScalarString
        # - test string decryption
    

    def test_walk_list_and_decrypt(self):
        """Walk list and decrypt its values."""
        # - test dict
        # - test list
        # - test ScalarString
        # - test string decryption

    def test_decrypt(self):
        """Test decrypt return a decrypted value."""

    # Test encryption
    def test_walk_and_encrypt(self):
        """Walk the branch recursively and encrypts its leaves."""
        # - test stash value
        # - test dict encryption
        # - test list values encryption
        # - test ScalarString
        # - test string encryption

    def test_walk_list_and_encrypt(self):
        """Walk a list contained in a branch and encrypts its values."""
        # - test stash value
        # - test dict encryption
        # - test list values encryption
        # - test ScalarString
        # - test string encryption

    def test_encrypt(self):        
        """Test encrypt return a encrypted value."""

    # Test keys management
    def test_get_key(self):
        """Test we obtain a 256 bits symetric key."""
        # - Test KMS key loading
        # - Test PGP key loading
        # - Test new key generation

    def test_get_key_from_kms(self):
        """Test we get the key form the KMS tree leave."""

    def test_encrypt_key_with_kms(self):
        """Test KMS encryption."""

    def test_get_key_from_pgp(self):
        """Test we get the key form the PGP tree leave."""

    def test_encrypt_key_with_pgp(self):
        """Test PGP encryption."""

    # Write file
    def test_write_file(self):
        """Test we can write a correct file with correct encoding."""

    # Open editor
    def test_run_editor(self):
        """Test we can run the editor with the specified file path."""

    # Panic errors
    def test_panic_writes_to_stderr(self):
        with mock.patch.object(builtins, 'print') as print_mock:
            with mock.patch("sys.exit") as sys_exit_mock:
                sops.panic("Foobar")
                print_mock.assert_called_with("Foobar", file=sys.stderr)
                sys_exit_mock.assert_called_with(1)

    def test_panic_handles_exit_error_code(self):
        with mock.patch.object(builtins, 'print'):
            with mock.patch("sys.exit") as sys_exit_mock:
                sops.panic("Foobar", 111)
                sys_exit_mock.assert_called_with(111)
