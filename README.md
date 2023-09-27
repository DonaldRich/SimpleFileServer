# SimpleFileServer
A simple file server for serving files from a folder on localhost.

SimpleFileServer was created to address the need for serving local files for development.

The fileserver is a single executable with the following features:
- The port used is specified as an argument.
- The root folder is specified as an argument. Only files relative to the specified root folder will be accessable.
- The response content type is set based on the file extensions. Common MimeTypes for HTML, Javascript, CSS, JSON, text, and image files are predefined.
- Additional MimeTypes can be added by supplying a JSON file (mime-types.json) containing the mapping from extensions to MimeTypes.
- By default only URLs to localhost or [::1] are allowed. If a clients.json file containing an array of clients is found these clients will be added to the allowed list. If this file contains only a single client of "*" client filtering will be turned off.

Use:
  SFS \<port\> \<root\>
