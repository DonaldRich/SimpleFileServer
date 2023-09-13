# SimpleFileServer
A simple file server for serving files from a folder on localhost.

SimpleFileServer was created to address the need for serving local files for development and local use without being blocked by CORS.

The fileserver is a single executable with the following features:
  The port used is specified as an argument.
  The root folder is specified as an argument. Only files relative to the specified root folder will be accessable.
  Common MimeTypes are set based on the file extensions.
  Uncommon MimeTypes can be added by supplying a JSON file containing the mapping from extensions to MimeTypes.

Use:
  SFS <port> <root>
