VIRTUALENV = virtualenv
VENV := $(shell echo $${VIRTUAL_ENV-.venv})
PYTHON = $(VENV)/bin/python
DEV_STAMP = $(VENV)/.dev_env_installed.stamp
INSTALL_STAMP = $(VENV)/.install.stamp

.IGNORE: clean
.PHONY: all install dev-requirements


all: build

build:
	python setup.py build

install: $(INSTALL_STAMP)

$(INSTALL_STAMP): $(PYTHON) setup.py
	$(VENV)/bin/pip install -U pip
	$(VENV)/bin/pip install -Ue .
	touch $(INSTALL_STAMP)

dev-requirements: $(INSTALL_STAMP) $(DEV_STAMP)
$(DEV_STAMP): $(PYTHON) dev-requirements.txt
	$(VENV)/bin/pip install tox
	$(VENV)/bin/pip install -Ur dev-requirements.txt
	touch $(DEV_STAMP)

virtualenv: $(PYTHON)
$(PYTHON):
	virtualenv $(VENV)

tests: dev-requirements
	$(VENV)/bin/tox

tests-once: install dev-requirements
	$(VENV)/bin/py.test --cov-report term-missing --cov sops tests/

functional-tests:
	gpg --import tests/sops_functional_tests_key.asc 2>&1 1>/dev/null || exit 0
	for type in yaml json txt; do \
		for ver in 2.6 2.7 3.4; do \
			echo "Testing Python$$ver $$type decryption" && \
			python$$ver sops/__init__.py -d example.$$type > /tmp/testdata.$$type && \
			echo "Testing Python$$ver $$type encryption" && \
			python$$ver sops/__init__.py -e -p "1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A" /tmp/testdata.$$type > /tmp/testdata$$ver.$$type; \
		done && \
		echo "Testing Python2.6 decryption of a 2.7 $$type file" && \
		python2.6 sops/__init__.py -d /tmp/testdata2.7.$$type > /dev/null && \
		echo "Testing Python2.6 decryption of a 3.4 $$type file" && \
		python2.6 sops/__init__.py -d /tmp/testdata3.4.$$type > /dev/null && \
		echo "Testing Python2.7 decryption of a 2.6 $$type file" && \
		python2.7 sops/__init__.py -d /tmp/testdata2.6.$$type > /dev/null && \
		echo "Testing Python2.7 decryption of a 3.4 $$type file" && \
		python2.7 sops/__init__.py -d /tmp/testdata3.4.$$type > /dev/null && \
		echo "Testing Python3.4 decryption of a 2.6 $$type file" && \
		python3.4 sops/__init__.py -d /tmp/testdata2.6.$$type > /dev/null && \
		echo "Testing Python3.4 decryption of a 2.7 $$type file" && \
		python3.4 sops/__init__.py -d /tmp/testdata2.7.$$type > /dev/null || exit 1; \
	done && \
	for ver in 2.6 2.7 3.4; do \
	done

functional-tests-once:
	gpg --import tests/sops_functional_tests_key.asc 2>&1 1>/dev/null || exit 0
	for type in yaml json txt; do \
		echo "Testing $$type decryption"; \
		python sops/__init__.py -d example.$$type > /tmp/testdata.$$type; \
		echo "Testing $$type encryption" ; \
		python sops/__init__.py -e -p "1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A" /tmp/testdata.$$type > /tmp/testdataenc.$$type; \
		echo "Testing $$type re-decryption" ; \
		python sops/__init__.py -d /tmp/testdataenc.$$type > /dev/null ; \
		echo "Testing removing PGP key to $$type encrypted file" ; \
		python sops/__init__.py -r --rm-pgp 85D77543B3D624B63CEA9E6DBC17301B491B3F21 /tmp/testdataenc.$$type ; \
	done
	echo "Testing round-trip on binary file"
	dd if=/dev/urandom of=/tmp/testdata-randomfile bs=1024 count=1024 2>&1 1>/dev/null
	python sops/__init__.py -e -p "1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A" /tmp/testdata-randomfile > /tmp/testdata-randomfile.enc
	python sops/__init__.py -d /tmp/testdata-randomfile.enc > /tmp/testdata-randomfile.dec
	if [ $$(shasum -a 256 /tmp/testdata-randomfile | cut -d ' ' -f 1) != $$(shasum -a 256 /tmp/testdata-randomfile.dec | cut -d ' ' -f 1) ]; then \
		echo "Binary file roundtrip failed, checksum doesn't match"; exit 0; \
	else \
		echo "Binary file roundtrip succeeded"; \
	fi;

pypi:
	$(PYTHON) setup.py sdist check upload --sign

clean:
	rm -rf *.pyc sops/*.pyc
	rm -rf __pycache__ sops/__pycache__
	rm -rf build/ dist/
	rm -fr .tox/ .venv/
	rm -fr .coverage
