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

try:
    from collections import OrderedDict
except ImportError:
    from ordereddict import OrderedDict

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

    def test_sops_branch_is_restored(self):
        m = mock.mock_open(read_data=sops.DEFAULT_YAML)
        b = {'kms': [ { 'arn': 'test' } ] }
        with mock.patch.object(builtins, 'open', m):
            tree = sops.load_file_into_tree('path', 'yaml',
                                            restore_sops=b)
            assert tree['sops']['kms'][0]['arn'] == 'test'

    def test_detect_filetype_handle_json(self):
        assert sops.detect_filetype("file.json") == "json"

    def test_detect_filetype_handle_yml(self):
        assert sops.detect_filetype("file.yml") == "yaml"

    def test_detect_filetype_handle_yaml(self):
        assert sops.detect_filetype("file.yaml") == "yaml"

    def test_detect_filetype_returns_text_if_unknown(self):
        assert sops.detect_filetype("file.xml") == "bytes"

    def test_verify_or_create_sops_branch(self):
        """Verify or create the sops branch"""
        # - sops is created if missing from tree
        # - kms arn is used
        # - pgp fp is used
        # - SOPS_KMS_ARN env variable is used
        # - SOPS_PGP_FP env variable is used
        # - panic error is raise and program quit with code 111 if
        #   nothing is defined
        log = logging.getLogger( "TreeTest.test_verify_or_create_sops_branch" )
        kms_arns = "arn:aws:kms:us-east-1:656532927350:key/" + \
                "920aff2e-c5f1-4040-943a-047fa387b27e+arn:aws:iam::" +\
                "927034868273:role/sops-dev, arn:aws:kms:ap-southeast-1:" + \
                "656532927350:key/9006a8aa-0fa6-4c14-930e-a2dfb916de1d"
        pgp_fps = "85D77543B3D624B63CEA9E6DBC17301B491B3F21," + \
                "C9CAB0AF1165060DB58D6D6B2653B624D620786D"
        tree = OrderedDict()
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
        key = os.urandom(32)
        tree = OrderedDict()
        with mock.patch.object(builtins, 'open', m):
            tree = sops.load_file_into_tree('path', 'yaml')
        tree['sops'] = dict()
        crypttree = sops.walk_and_encrypt(tree, key)
        assert crypttree['example_key'].startswith("ENC[AES256_GCM,data:")
        assert isinstance(crypttree['example_array'], list)
        assert len(crypttree['example_array']) == 2

    def test_walk_and_encrypt_unencrypted(self):
        """Walk the branch with whitelisted leaves and verify they are untouched."""
        m = mock.mock_open(read_data="""{
            "str_unencrypted": "STRING",
            "list_unencrypted": ["A", "B"],
            "dict_unencrypted": {
                "key": "value"
            },
            "foo": "bar"
        }""")
        # Verify data stays unencrypted upon encryption
        key = os.urandom(32)
        tree = OrderedDict()
        with mock.patch.object(builtins, 'open', m):
            tree = sops.load_file_into_tree('path', 'json')
        tree['sops'] = dict()
        crypttree = sops.walk_and_encrypt(OrderedDict(tree), key, isRoot=True)
        assert crypttree['str_unencrypted'] == 'STRING'
        assert crypttree['list_unencrypted'] == ['A', 'B']
        assert crypttree['dict_unencrypted'] == {"key": "value"}
        assert crypttree['foo'].startswith("ENC[AES256_GCM,data:")

        # Verify we have a MAC and it includes unencrypted values
        assert tree['sops']['mac'].startswith("ENC[AES256_GCM,data:")
        empty_tree = OrderedDict()
        empty_tree['sops'] = dict()
        sops.walk_and_encrypt(OrderedDict(empty_tree), key, isRoot=True)
        tree_mac = sops.decrypt(tree['sops']['mac'], key,
                                aad=tree['sops']['lastmodified'].encode('utf-8'))
        empty_tree_mac = sops.decrypt(empty_tree['sops']['mac'], key,
                                      aad=empty_tree['sops']['lastmodified'].encode('utf-8'))
        assert tree_mac != empty_tree_mac

    def test_walk_and_encrypt_and_decrypt(self):
        """Test a roundtrip on the tree encryption/decryption code"""
        m = mock.mock_open(read_data=sops.DEFAULT_JSON)
        key = os.urandom(32)
        tree = OrderedDict()
        with mock.patch.object(builtins, 'open', m):
            tree = sops.load_file_into_tree('path', 'json')
        tree['sops'] = dict()
        crypttree = sops.walk_and_encrypt(OrderedDict(tree), key, isRoot=True)
        cleartree = sops.walk_and_decrypt(OrderedDict(crypttree), key, isRoot=True)
        assert cleartree == tree

    def test_numbers_encrypt_and_decrypt(self):
        """Test encryption/decryption of numbers"""
        m = mock.mock_open(read_data='{"a":1234,"b":[567,890.123],"c":5.4999517527e+10}')
        key = os.urandom(32)
        tree = OrderedDict()
        with mock.patch.object(builtins, 'open', m):
            tree = sops.load_file_into_tree('path', 'json')
        tree['sops'] = dict()
        crypttree = sops.walk_and_encrypt(OrderedDict(tree), key, isRoot=True)
        assert tree['sops']['mac'].startswith("ENC[AES256_GCM,data:")
        cleartree = sops.walk_and_decrypt(OrderedDict(crypttree), key, isRoot=True)
        assert cleartree == tree

    def test_bytes_encrypt_and_decrypt(self):
        """Test encryption/decryption of numbers"""
        key = os.urandom(32)
        tree = OrderedDict()
        tree['data'] = os.urandom(4096)
        tree['sops'] = dict()
        crypttree = sops.walk_and_encrypt(OrderedDict(tree), key, isRoot=True)
        assert tree['sops']['mac'].startswith("ENC[AES256_GCM,data:")
        cleartree = sops.walk_and_decrypt(OrderedDict(crypttree), key, isRoot=True)
        assert cleartree == tree

    def test_walk_and_encrypt_and_decrypt_unencrypted(self):
        """Test a roundtrip on the tree encryption/decryption code with unencrypted leaves"""
        m = mock.mock_open(read_data="""{
            "str_unencrypted": "STRING",
            "list_unencrypted": ["A", "B"],
            "dict_unencrypted": {
                "key": "value"
            },
            "foo": "bar"
        }""")
        key = os.urandom(32)
        tree = OrderedDict()
        with mock.patch.object(builtins, 'open', m):
            tree = sops.load_file_into_tree('path', 'json')
        tree['sops'] = dict()
        crypttree = sops.walk_and_encrypt(OrderedDict(tree), key, isRoot=True)
        assert tree['sops']['mac'].startswith("ENC[AES256_GCM,data:")
        cleartree = sops.walk_and_decrypt(OrderedDict(crypttree), key, isRoot=True)
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
        origin = "AAAAAAAA"
        key = os.urandom(32)
        aad = os.urandom(32)
        clearstr = sops.decrypt(sops.encrypt(origin, key, aad=aad), key, aad=aad)
        assert clearstr == origin

    def test_add_kms_master_keys(self):
        """ test adding a kms master key to an existing tree """
        tree = {'sops': { 'kms': [ {'arn': 'arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e' } ] } }
        newkms = 'arn:aws:kms:us-east-1:656532927350:key/920abb2e-c2b3-9090-943a-047fa387f3ac+arn:aws:iam::927034868273:role/sops-dev-xyz'
        assert len(tree['sops']['kms']) == 1
        tree = sops.add_new_master_keys(tree, newkms, '')
        assert tree['sops']['kms'][0]['arn'] == 'arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e'
        assert tree['sops']['kms'][1]['arn'] == 'arn:aws:kms:us-east-1:656532927350:key/920abb2e-c2b3-9090-943a-047fa387f3ac'
        assert tree['sops']['kms'][1]['role'] == 'arn:aws:iam::927034868273:role/sops-dev-xyz'
                            
    def test_add_pgp_master_keys_where_none_existed(self):
        """ test adding a pgp master key to an existing tree
            that does not have any pgp master key yet
        """
        tree = {'sops': { 'kms': [ {'arn': 'arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e' } ] } }
        newpgp = 'E60892BB9BD89A69F759A1A0A3D652173B763E8F'
        tree = sops.add_new_master_keys(tree, '', newpgp)
        assert tree['sops']['kms'][0]['arn'] == 'arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e'
        assert tree['sops']['pgp'][0]['fp'] == 'E60892BB9BD89A69F759A1A0A3D652173B763E8F'
                            
    def test_add_pgp_master_keys(self):
        """ test adding a pgp master key to an existing tree """
        tree = {'sops': { 'pgp': [ {'fp': '1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A' } ] } }
        newpgp = 'E60892BB9BD89A69F759A1A0A3D652173B763E8F'
        assert len(tree['sops']['pgp']) == 1
        tree = sops.add_new_master_keys(tree, '', newpgp)
        assert tree['sops']['pgp'][0]['fp'] == '1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A'
        assert tree['sops']['pgp'][1]['fp'] == 'E60892BB9BD89A69F759A1A0A3D652173B763E8F'
                            
    def test_add_kms_master_keys_where_none_existed(self):
        """ test adding a kms master key to an existing tree
            that does not have any kms master key yet
        """
        tree = {'sops': { 'pgp': [ {'fp': '1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A' } ] } }
        newkms = 'arn:aws:kms:us-east-1:656532927350:key/920abb2e-c2b3-9090-943a-047fa387f3ac+arn:aws:iam::927034868273:role/sops-dev-xyz'
        tree = sops.add_new_master_keys(tree, newkms, '')
        assert tree['sops']['pgp'][0]['fp'] == '1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A'
        assert tree['sops']['kms'][0]['arn'] == 'arn:aws:kms:us-east-1:656532927350:key/920abb2e-c2b3-9090-943a-047fa387f3ac'
        assert tree['sops']['kms'][0]['role'] == 'arn:aws:iam::927034868273:role/sops-dev-xyz'

    def test_rm_kms_master_keys(self):
        """ test removing a kms master key to an existing tree """
        tree = {'sops': { 'kms': [
                    {'arn': 'arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e' },
                    {'arn': 'arn:aws:kms:us-east-1:656532927350:key/920abb2e-c2b3-9090-943a-047fa387f3ac' }
                ] } }
        rmkms = 'arn:aws:kms:us-east-1:656532927350:key/920abb2e-c2b3-9090-943a-047fa387f3ac+arn:aws:iam::927034868273:role/sops-dev-xyz'
        assert len(tree['sops']['kms']) == 2
        tree = sops.remove_master_keys(tree, rmkms, '')
        assert tree['sops']['kms'][0]['arn'] == 'arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e'
        assert len(tree['sops']['kms']) == 1
                            
    def test_rm_pgp_master_keys_where_none_existed(self):
        """ test removing a pgp master key to an existing tree
            that does not have any pgp master key yet
        """
        tree = {'sops': { 'kms': [ {'arn': 'arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e' } ] } }
        rmpgp = 'E60892BB9BD89A69F759A1A0A3D652173B763E8F'
        tree = sops.remove_master_keys(tree, '', rmpgp)
        assert tree['sops']['kms'][0]['arn'] == 'arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e'
        assert len(tree['sops']['kms']) == 1
        assert 'pgp' not in tree['sops']
 
    def test_rm_pgp_master_keys(self):
        """ test removing a pgp master key to an existing tree """
        tree = {'sops': { 'pgp': [
                    {'fp': '1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A' },
                    {'fp': 'E60892BB9BD89A69F759A1A0A3D652173B763E8F' }
                ] } }
        rmpgp = 'E60892BB9BD89A69F759A1A0A3D652173B763E8F'
        assert len(tree['sops']['pgp']) == 2
        tree = sops.remove_master_keys(tree, '', rmpgp)
        assert len(tree['sops']['pgp']) == 1
        assert tree['sops']['pgp'][0]['fp'] == '1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A'
                            
    def test_rm_kms_master_keys_where_none_existed(self):
        """ test removing a kms master key to an existing tree
            that does not have any kms master key yet
        """
        tree = {'sops': { 'pgp': [ {'fp': '1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A' } ] } }
        rmkms = 'arn:aws:kms:us-east-1:656532927350:key/920abb2e-c2b3-9090-943a-047fa387f3ac+arn:aws:iam::927034868273:role/sops-dev-xyz'
        tree = sops.remove_master_keys(tree, rmkms, '')
        assert tree['sops']['pgp'][0]['fp'] == '1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A'
        assert len(tree['sops']['pgp']) == 1
        assert 'kms' not in tree['sops']
                            
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

    def test_valid_json_syntax(self):
        m = mock.mock_open(read_data=sops.DEFAULT_JSON)
        with mock.patch.object(builtins, 'open', m):
            assert sops.validate_syntax('path', 'json') == True

    def test_invalid_json_syntax(self):
        m = mock.mock_open(read_data='{,,,,,}')
        with mock.patch.object(builtins, 'open', m):
            with self.assertRaises(ValueError):
                sops.validate_syntax('path', 'json')

    def test_valid_yaml_syntax(self):
        m = mock.mock_open(read_data=sops.DEFAULT_YAML)
        with mock.patch.object(builtins, 'open', m):
            assert sops.validate_syntax('path', 'yaml') == True

    def test_bytes_syntax(self):
        m = mock.mock_open(read_data=sops.DEFAULT_TEXT)
        with mock.patch.object(builtins, 'open', m):
            assert sops.validate_syntax('path', 'bytes') == True

    def test_subtree(self):
        """Extract a subtree from a document."""
        m = mock.mock_open(read_data=sops.DEFAULT_YAML)
        key = os.urandom(32)
        tree = OrderedDict()
        with mock.patch.object(builtins, 'open', m):
            tree = sops.load_file_into_tree('path', 'yaml')
        ntree = sops.truncate_tree(dict(tree), '["example"]["nested"]["values"]')
        assert ntree == tree["example"]["nested"]["values"]
        ntree = sops.truncate_tree(dict(tree), '["example_array"][1]')
        assert ntree == tree["example_array"][1]

    def test_version_comparison(self):
        assert sops.A_is_newer_than_B('1.9', '1.8') == True
        assert sops.A_is_newer_than_B('2.1', '1.8') == True
        assert sops.A_is_newer_than_B('1.7', '1.8') == False
        assert sops.A_is_newer_than_B('1.8.2', '1.8') == True
        assert sops.A_is_newer_than_B('1.7.9', '1.8') == False
        assert sops.A_is_newer_than_B('1.7', '1.8.1') == False
