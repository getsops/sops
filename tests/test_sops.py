# -*- coding: utf-8 -*-
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
#
# Contributor: Julien Vehent <jvehent@mozilla.com> [:ulfr]
# Contributor: Daniel Thornton <daniel@relud.com>
# Contributor: Alexis Metaireau <alexis@mozilla.com> [:alexis]
# Contributor: Rémy Hubscher <natim@mozilla.com> [:natim]

import logging
import unittest2
import mock
import os
import sys

import sops

if sys.version_info[0] == 2:
    import __builtin__ as builtins
else:
    import builtins


class TreeTest(unittest2.TestCase):

    def test_json_loader_is_used_on_json_filetype(self):
        m = mock.mock_open(read_data=sops.DEFAULT_JSON)
        with mock.patch.object(builtins, 'open', m):
            tree = sops.load_file_into_tree('path', 'json')
            assert tree['example_key'] == 'example_value'

    def test_yaml_loader_is_used_on_yaml_filetype(self):
        m = mock.mock_open(read_data=sops.DEFAULT_YAML)
        with mock.patch.object(builtins, 'open', m):
            tree = sops.load_file_into_tree('path', 'yaml')
            assert tree['example_key'] == 'example_value'

    #def test_text_loader_is_used_on_text_filetype(self):
    #    m = mock.mock_open(read_data=sops.DEFAULT_TEXT)
    #    with mock.patch.object(builtins, 'open', m):
    #        tree = sops.load_file_into_tree('path', 'text')
    #        assert tree['data'].startswith(sops.DEFAULT_TEXT[0:15])

    @mock.patch('sops.json.load')
    def test_example_with_a_mocked_call(self, json_mock):
        m = mock.mock_open(read_data='"content"')
        with mock.patch.object(builtins, 'open', m):
            sops.load_file_into_tree('path', 'json')
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
        log= logging.getLogger( "TreeTest.test_verify_or_create_sops_branch" )
        kms_arns = "arn:aws:kms:us-east-1:656532927350:key/" + \
                "920aff2e-c5f1-4040-943a-047fa387b27e+arn:aws:iam::" +\
                "927034868273:role/sops-dev, arn:aws:kms:ap-southeast-1:" + \
                "656532927350:key/9006a8aa-0fa6-4c14-930e-a2dfb916de1d"
        pgp_fps = "85D77543B3D624B63CEA9E6DBC17301B491B3F21," + \
                "C9CAB0AF1165060DB58D6D6B2653B624D620786D"
        tree = dict()
        tree, ign = sops.verify_or_create_sops_branch(tree,
                                                      kms_arns=kms_arns,
                                                      pgp_fps=pgp_fps)
        log.debug("%s", tree)
        assert len(tree['sops']['kms']) == 2
        assert tree['sops']['kms'][0]['arn'] == "arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e"
        assert tree['sops']['kms'][0]['role'] == "arn:aws:iam::927034868273:role/sops-dev"
        assert tree['sops']['kms'][1]['arn'] == "arn:aws:kms:ap-southeast-1:656532927350:key/9006a8aa-0fa6-4c14-930e-a2dfb916de1d"
        assert len(tree['sops']['pgp']) == 2
        assert tree['sops']['pgp'][0]['fp'] == "85D77543B3D624B63CEA9E6DBC17301B491B3F21"
        assert tree['sops']['pgp'][1]['fp'] == "C9CAB0AF1165060DB58D6D6B2653B624D620786D"

    def test_update_sops_branch(self):
        """ If master keys have been added to the SOPS branch, encrypt the data key
            with them, and store the new encrypted values.
        """
        # - verify data key gets encrypted with new master key

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

    # Test encryption
    def test_walk_and_encrypt(self):
        """Walk the branch recursively and encrypts its leaves."""
        # - test dict encryption
        # - test list values encryption
        # - test ScalarString
        # - test string encryption
        # TODO: 
        # - test stash value
        m = mock.mock_open(read_data=sops.DEFAULT_YAML)
        tree = dict()
        key = os.urandom(32)
        with mock.patch.object(builtins, 'open', m):
            tree = sops.load_file_into_tree('path', 'yaml')
        crypttree = sops.walk_and_encrypt(tree, key)
        assert crypttree['example_key'].startswith("ENC[AES256_GCM,data:")
        assert isinstance(crypttree['example_array'], list)
        assert len(crypttree['example_array']) == 2

    def test_walk_and_encrypt_and_decrypt(self):
        """Test a roundtrip on the tree encryption/decryption code"""
        m = mock.mock_open(read_data=sops.DEFAULT_JSON)
        tree = dict()
        key = os.urandom(32)
        with mock.patch.object(builtins, 'open', m):
            tree = sops.load_file_into_tree('path', 'json')
        crypttree = sops.walk_and_encrypt(tree, key)
        cleartree = sops.walk_and_decrypt(crypttree, key)
        assert cleartree == tree

    def test_walk_list_and_encrypt(self):
        """Walk a list contained in a branch and encrypts its values."""
        # - test stash value
        # - test dict encryption
        # - test list values encryption
        # - test ScalarString
        # - test string encryption

    def test_encrypt(self):
        """Test encrypt return a encrypted value."""
        cryptstr = sops.encrypt("AAAAAAA", os.urandom(32))
        assert cryptstr.startswith("ENC[AES256_GCM,data:")
        assert cryptstr[-1:] == "]"

    def test_encrypt_decrypt(self):
        """Test a roundtrip in the encryption/decryption code"""
        key = os.urandom(32)
        cryptstr = sops.encrypt("AAAAAAA", key)
        clearstr = sops.decrypt(cryptstr, key)
        assert clearstr == "AAAAAAA"

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
                print_mock.assert_called_with("PANIC: Foobar", file=sys.stderr)
                sys_exit_mock.assert_called_with(1)

    def test_panic_handles_exit_error_code(self):
        with mock.patch.object(builtins, 'print'):
            with mock.patch("sys.exit") as sys_exit_mock:
                sops.panic("Foobar", 111)
                sys_exit_mock.assert_called_with(111)
