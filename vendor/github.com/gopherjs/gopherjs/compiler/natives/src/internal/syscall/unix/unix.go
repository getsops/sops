// +build js

package unix

const randomTrap = 0
const fstatatTrap = 0

func IsNonblock(fd int) (nonblocking bool, err error) {
	return false, nil
}
