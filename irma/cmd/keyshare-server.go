package cmd

import (
	"github.com/go-errors/errors"
	irma "github.com/privacybydesign/irmago"
	"github.com/privacybydesign/irmago/internal/keysharecore"
	"github.com/privacybydesign/irmago/server"
	"github.com/privacybydesign/irmago/server/keyshare/keyshareserver"
	"github.com/sietseringers/cobra"
	"github.com/sietseringers/viper"
)

var keyshareServerCmd = &cobra.Command{
	Use:   "server",
	Short: "IRMA keyshare server",
	Run: func(command *cobra.Command, args []string) {
		conf, err := configureKeyshareServer(command)
		if err != nil {
			die("failed to read configuration", err)
		}

		// Create main server
		keyshareServer, err := keyshareserver.New(conf)
		if err != nil {
			die("", err)
		}

		runServer(keyshareServer, conf.Logger)
	},
}

func init() {
	keyshareRootCmd.AddCommand(keyshareServerCmd)

	flags := keyshareServerCmd.Flags()
	flags.SortFlags = false
	flags.StringP("config", "c", "", "path to configuration file")
	flags.StringP("schemes-path", "s", irma.DefaultSchemesPath(), "path to irma_configuration")
	flags.String("schemes-assets-path", "", "if specified, copy schemes from here into --schemes-path")
	flags.Int("schemes-update", 60, "update IRMA schemes every x minutes (0 to disable)")
	flags.StringP("privkeys", "k", "", "path to IRMA private keys")
	flags.StringP("url", "u", "", "external URL to server to which the IRMA client connects, \":port\" being replaced by --port value")

	flags.IntP("port", "p", 8080, "port at which to listen")
	flags.StringP("listen-addr", "l", "", "address at which to listen (default 0.0.0.0)")
	flags.Lookup("port").Header = `Server address and port to listen on`

	flags.String("db-type", string(keyshareserver.DBTypePostgres), "Type of database to connect keyshare server to")
	flags.String("db", "", "Database server connection string")
	flags.Lookup("db-type").Header = `Database configuration`

	flags.String("jwt-privkey", "", "Private jwt key of keyshare server")
	flags.String("jwt-privkey-file", "", "Path to file containing private jwt key of keyshare server")
	flags.Int("jwt-privkey-id", 0, "Key identifier of keyshare server public key matching used private key")
	flags.String("jwt-issuer", keysharecore.JWTIssuerDefault, "JWT issuer used in \"iss\" field")
	flags.Int("jwt-pin-expiry", keysharecore.JWTPinExpiryDefault, "Expiry of PIN JWT in seconds")
	flags.String("storage-primary-keyfile", "", "Primary key used for encrypting and decrypting secure containers")
	flags.StringSlice("storage-fallback-keyfile", nil, "Fallback key(s) used to decrypt older secure containers")
	flags.Lookup("jwt-privkey").Header = `Cryptographic keys`

	flags.String("keyshare-attribute", "", "Attribute identifier that contains username")
	flags.Lookup("keyshare-attribute").Header = `Keyshare server attribute issued during registration`

	flags.String("email-server", "", "Email server to use for sending email address confirmation emails")
	flags.String("email-hostname", "", "Hostname used in email server tls certificate (leave empty when mail server does not use tls)")
	flags.String("email-username", "", "Username to use when authenticating with email server")
	flags.String("email-password", "", "Password to use when authenticating with email server")
	flags.String("email-from", "", "Email address to use as sender address")
	flags.String("default-language", "en", "Default language, used as fallback when users preferred language is not available")
	flags.StringToString("registration-email-subjects", nil, "Translated subject lines for the registration email")
	flags.StringToString("registration-email-files", nil, "Translated emails for the registration email")
	flags.StringToString("verification-url", nil, "Base URL for the email verification link (localized)")
	flags.Lookup("email-server").Header = `Email configuration (leave empty to disable sending emails)`

	flags.String("tls-cert", "", "TLS certificate (chain)")
	flags.String("tls-cert-file", "", "path to TLS certificate (chain)")
	flags.String("tls-privkey", "", "TLS private key")
	flags.String("tls-privkey-file", "", "path to TLS private key")
	flags.Bool("no-tls", false, "Disable TLS")
	flags.Lookup("tls-cert").Header = `TLS configuration (leave empty to disable TLS)`

	flags.CountP("verbose", "v", "verbose (repeatable)")
	flags.BoolP("quiet", "q", false, "quiet")
	flags.Bool("log-json", false, "Log in JSON format")
	flags.Bool("production", false, "Production mode")
	flags.Lookup("verbose").Header = `Other options`
}

func configureKeyshareServer(cmd *cobra.Command) (*keyshareserver.Configuration, error) {
	readConfig(cmd, "keyshareserver", "keyshareserver", []string{".", "/etc/keyshareserver"}, nil)

	// And build the configuration
	conf := &keyshareserver.Configuration{
		Configuration:      configureIRMAServer(),
		EmailConfiguration: configureEmail(),

		DBType:    keyshareserver.DBType(viper.GetString("db-type")),
		DBConnStr: viper.GetString("db-str"),

		JwtKeyID:                viper.GetUint32("jwt-privkey-id"),
		JwtPrivateKey:           viper.GetString("jwt-privkey"),
		JwtPrivateKeyFile:       viper.GetString("jwt-privkey-file"),
		JwtIssuer:               viper.GetString("jwt-issuer"),
		JwtPinExpiry:            viper.GetInt("jwt-pin-expiry"),
		StoragePrimaryKeyFile:   viper.GetString("storage-primary-keyfile"),
		StorageFallbackKeyFiles: viper.GetStringSlice("storage-fallback-keyfile"),

		KeyshareAttribute: irma.NewAttributeTypeIdentifier(viper.GetString("keyshare-attribute")),

		RegistrationEmailSubjects: viper.GetStringMapString("registration-email-subjects"),
		RegistrationEmailFiles:    viper.GetStringMapString("registration-email-files"),
		VerificationURL:           viper.GetStringMapString("verification-url"),
	}

	if conf.Production && conf.DBType != keyshareserver.DBTypePostgres {
		return nil, errors.New("in production mode, db-type must be postgres")
	}

	conf.URL = server.ReplacePortString(viper.GetString("url"), viper.GetInt("port"))

	return conf, nil
}
