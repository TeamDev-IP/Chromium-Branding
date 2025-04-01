# Chromium Branding

This repository provides a command line tool for applying custom branding to the Chromium binaries.

[JxBrowser](https://teamdev.com/jxbrowser) and [DotNetBrowser](https://teamdev.com/dotnetbrowser) libraries use their own Chromium builds with the default branding. Chromium is running in a separate process. The process name is `Chromium` and it has the default icon, version, copyright, etc. For a better user experience, the software that uses JxBrowser or DotNetBrowser may want to customize the process name, icon, version, etc. of the Chromium binaries.

This command line tool allows you to customize the Chromium binaries by changing the process name, icon, version, etc. The tool supports Windows, macOS, and Linux.

## Prerequisites

Make sure you install [Go](https://go.dev/dl/) 1.20+ and add it to the `PATH` environment variable.

### Windows

On Windows, you also have to install the [up-to-date Windows SDK](https://developer.microsoft.com/en-us/windows/downloads/windows-sdk/).

**Important**: if you install or have already installed it to the **custom location**, you have to add its bin directory for the
current platform architecture to PATH.
For example, for default location, it will be `C:\Program Files (x86)\Windows Kits\10\bin\10.0.26100.0\x64`.

## Building

Run the following command in the root directory of the repository:

### Windows

```sh
go build -o chromium_branding.exe
```

### macOS/Linux

```sh
go build -o chromium_branding
```

You will find the `chromium_branding(.exe)` executable there.

## Usage

Take [params.json](params.json) and modify it with the necessary branding information or create a new one. The JSON file should have the following structure:

```JSON
{
  "version": "1.2.3",
  "win": {
    "executableName": "myapp",
    "processDisplayName": "My App",
    "legalCopyright": "© 2025 MyCompany",
    "author": "Me",
    "productName": "MyApp",
    "icoPath": "assets/app.ico",
    "signCommand": "echo @@BINARY_PATH@@"
  },
  "mac": {
    "bundle": {
      "name": "MyApp",
      "id": "com.mycompany.myapp"
    },
    "icnsPath": "assets/app.icns",
    "codesignIdentity": "${CODESIGN_IDENTITY}",
    "codesignEntitlements": "assets/entitlements.plist",
    "teamID": "${TEAM_ID}",
    "appleID": "${APPLE_ID}",
    "password": "${PASSWORD}"
  },
  "linux": {
    "executableName": "myapp"
  }
}
```

Here's the description of the JSON parameters:

| Parameter                  | Description                                                                                                                                                 |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `version`                  | The version of the app.                                                                                                                                     |
| `win.executableName`       | The name of the Windows executable without the `.exe` extension.                                                                                            |
| `win.processDisplayName`   | The name that will be associated with the process in the `Processes` list in Task Manager.                                                                  |
| `win.legalCopyright`       | The legal copyright property of the executable file.                                                                                                        |
| `win.author`               | The author property of the executable file.                                                                                                                 |
| `win.productName`          | The legal product name property of the executable file.                                                                                                     |
| `win.icoPath`              | The path to the `.ico` file that represents the Windows app icon.                                                                                           |
| `win.signCommand`          | The command that will be used to sign the executable file.  the `@@BINARY_PATH@@"` placeholder will be replaced with the actual file path during execution. |
| `mac.bundle.name`          | The name of the resulting macOS app bundle.                                                                                                                 |
| `mac.bundle.id`            | The bundle ID that will be associated with the app.                                                                                                         |
| `mac.icnsPath`             | The path to the `.icns` file that represents the macOS app icon.                                                                                            |
| `mac.codesignIdentity`     | The identity that will be used to sign the macOS app bundle.                                                                                                |
| `mac.codesignEntitlements` | The path to the entitlements file that will be used to sign the macOS app bundle.                                                                           |
| `mac.teamID`               | The team ID that will be used to sign the macOS app bundle.                                                                                                 |
| `mac.appleID`              | The Apple ID that will be used to notarize the macOS app bundle.                                                                                            |
| `mac.password`             | The password for the Apple ID that will be used to notarize the macOS app bundle.                                                                           |
| `linux.executableName`     | The name of the branded executable on Linux.                                                                                                                |


Run the following command in the terminal to customize the Chromium binaries:

### Windows

```sh
chromium_branding.exe -p <params-json> -b <chromium-binaries-path> -o <output-dir>
```

### macOS/Linux

```sh
./chromium_branding -p <params-json> -b <chromium-binaries-path> -o <output-dir>
```

The customized Chromium binaries will be saved in the specified output directory.

**Important**: the tool will create a special `executable.name` file in the output directory. **Do not delete this file because it's necessary to run JxBrowser/DotNetBrowser with customized Chromium binaries.**

## Signing and notarizing

The original Chromium binaries deployed with JxBrowser and DotNetBrowser are signed with the TeamDev certificate and notarized by Apple. When you customize the Chromium binaries, you lose the original signature and notarization.

If you want to deploy the customized Chromium binaries with your software, you need to sign them with your own certificate and notarize.

The `win.signCommand` parameter in the `params.json` file allows you to sign the Windows executable.

The `mac.codesignIdentity`, `mac.codesignEntitlements`, `mac.teamID`, `mac.appleID`, and `mac.password` parameters in the `params.json` file allow you to sign and notarize the macOS app bundle.
