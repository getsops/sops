"""
Patch go.mod so that the lines 'go xxx' to 'toolchain xxx' are as in git's
HEAD.

This is necessary since newer 'go mod tidy' versions tend to modify these
lines. Since we check in CI that 'go mod tidy' does not change go.mod, this
causes CI to fail.
"""

import subprocess


def split_go_mod(contents: str) -> tuple[list[str], list[str], list[str]]:
    """
    Given the contents of go.mod, splits it into three lists of lines
    (with endings):
    1. The lines before 'go';
    2. The lines starting with 'go' and ending with 'toolchain';
    3. The lines after 'toolchain'.
    """
    parts: tuple[list[str], list[str], list[str]] = ([], [], [])
    index = 0
    for line in contents.splitlines(keepends=True):
        next_index = index
        if line.startswith('go '):
            index = next_index = 1
        if line.startswith('toolchain '):
            next_index = 2
        parts[index].append(line)
        index = next_index
    return parts


def get_file_contents_from_git_revision(filename: str, revision: str) -> str:
    """
    Get the file contents of ``filename`` from Git revision ``revision``.
    """
    p = subprocess.run(
        ['git', 'show', f'{revision}:{filename}'],
        stdout=subprocess.PIPE,
        check=True,
        encoding='utf-8',
    )
    return p.stdout


def read_file(filename: str) -> str:
    """
    Read the file's contents.
    """
    with open(filename, 'r', encoding='utf-8') as f:
        return f.read()


def write_file(filename: str, contents: str) -> None:
    """
    Write the file's contents.
    """
    with open(filename, 'w', encoding='utf-8') as f:
        f.write(contents)


def main():
    """
    Patches go.mod.
    """
    filename = 'go.mod'
    _, go_versions, __ = split_go_mod(
        get_file_contents_from_git_revision(filename, 'HEAD')
    )
    head, _, tail = split_go_mod(read_file(filename))
    lines = head + go_versions + tail
    write_file(filename, ''.join(lines))


if __name__ == '__main__':
    main()
