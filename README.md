# doc-generator-poc
This is a proof-of-concept on generating `.docx` document based on a template using Golang. The template is also in a 
`.docx` format. The app generate the document by unzipping the template document and search for `word/document.xml` file.
It will load the xml file and create new `template` based on that.

## Reference
* [.docx file format](https://wiki.fileformat.com/word-processing/docx/) 