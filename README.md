# Chromium Branding

This repository provides a command line tool for applying custom branding to Chromium binaries.

## Chromium binaries

We build a proprietary project called Platinum, which produces a set of native libraries and executables.
These libraries expose APIs used by our JxBrowser, DotNetBrowser, and Molybden projects.

We ship the appropriate version of the Chromium binaries with each specific release of JxBrowser, DotNetBrowser, or Molybden.
 
## Chromium process branding

The binaries produced by the Platinum project contain the default Chromium branding, which can be confusing for users of applications that rely on Chromium-based libraries.

This repository offers a tool to patch these binaries, providing a custom application icon, process name,
and other resources. This is essential for rebranding the native Chromium process to match the
specific application that uses our Chromium-based libraries.