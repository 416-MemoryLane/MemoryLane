# MemoryLane

## Table of Contents

- [Requirements](#requirements)
- [Usage](#usage)
- [Restrictions & Assumptions](#restrictions--assumptions)

### Requirements

- Go version 1.16 or later. You can download and install it from the official website: https://golang.org/dl/.
- Node.js version 14 or later. You can download and install it from the official website: https://nodejs.org/en/download/.
- Yarn package manager. You can install it using the following command: `npm install -g yarn`.

**Note**: Make sure that Go, Node.js, and Yarn are properly installed and configured before running the application.

### Usage

From the project root directory:

1. Modify the file permissions and allow execution of the start script `chmod u+x ./start`
2. Execute the bash script `./start <username> <password>`.

For windows users,

1. Execute the powershell script `./start.ps1 <username> <password>`

### Restrictions & Assumptions

- Starting the script with an unregistred user automatically creates a new user.
- We assume that there will be one user per repository instance. Once you have logged in successfully, you will not be able to change users.
  - To reset the state of the repository delete the `/.env` file and `/memory-lane-gallery` directory from the project root directory.
