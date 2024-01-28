#!/usr/bin/env python3

# Needs antsibull-changelog 0.24.0+ installed.

import sys

from antsibull_changelog.rendering.markdown import render_as_markdown


def main():
	with open('README.rst', 'rt', encoding='utf-8') as f:
		rst_data = f.read()

	result = render_as_markdown(rst_data, parser_name='restructuredtext', source_path='README.rst')

	if result.unsupported_class_names:
		print(f'Unknown class names: {sorted(result.unsupported_class_names)}', file=sys.stderr)

	with open('README.md', 'wt', encoding='utf-8') as f:
		f.write('<!-- THIS FILE HAS BEEN AUTOMATICALLY CONVERTED FROM README.rst.\n     DO NOT MODIFY THIS FILE, MODIFY README.rst INSTEAD! -->\n\n')
		f.write(result.output)
		f.write('\n')


main()
