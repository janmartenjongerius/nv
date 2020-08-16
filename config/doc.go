/*
Package config describes the structure of a Variable and lists variables for the
current Environment.

During the init-phase of the application, the Environment variable is populated
with the result of os.Environ.

Exporting variables
Using an Exporter, Variable structs can be written a provided io.Writer, for a
given format.

Encoding
Variables can be converted to a textual representation and back, using
Decoder and Encoder instances.

An Encoding is the combination of an Encoder and Decoder, which can be
registered against a given format.

Additional encoders can be added through the RegisterEncoding function.
This allows for separation of concerns internally, but also for plugins to
implement their own encoding.

The function NewEncoding allows to get the encoder that is registered against
the provided format. This function will return an error when no encoding is
known for the given format, even if the encoding is registered later on.

To process Variables for a given format, without being concerned about whether
the format exists, the function WithEncoding can be used. Since an Encoding may
be registered at runtime, this strategy allows to queue processing of Variables
before the Encoding is available. The moment the Encoding is registered, all
corresponding EncodingCallback functions are invoked. When the Encoding is
already present, this happens the moment WithEncoding is invoked.

When explicitly checking against an available Encoding, the functions
HasEncoding and GetEncodings will be of use.
*/
package config
