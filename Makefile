all:
	python setup.py build

install:
	python setup.py install

requirements:
	sudo pip install -U -r requirements.txt

dev-requirements:
	sudo pip install -U -r dev-requirements.txt

tests:
	py.test tests
test: tests

pypi:
	python setup.py sdist check upload --sign

clean:
	rm -rf *pyc
	rm -rf build
	rm -rf __pycache__
	rm -rf dist
