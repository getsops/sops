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

pypi:
	$(PYTHON) setup.py sdist check upload --sign

clean:
	rm -rf *.pyc sops/*.pyc
	rm -rf __pycache__ sops/__pycache__
	rm -rf build/ dist/
	rm -fr .tox/ .venv/
	rm -fr .coverage
