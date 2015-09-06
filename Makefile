all:
	./setup.py build

install:
	./setup.py install

rpm:
	fpm -s python -t rpm -d pytz -d python-requests-futures ./setup.py

deb:
	fpm -s python -t deb ./setup.py

tests: test
test:
	python ./sops -d example.yaml

pypi:
	 python setup.py sdist check upload --sign

clean:
	rm -rf *pyc
	rm -rf build
	rm -rf __pycache__
	rm -rf dist
