// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/spf13/cobra"
)

type hadolintExecuteOptions struct {
	ConfigurationURL      string `json:"configurationUrl,omitempty"`
	ConfigurationUsername string `json:"configurationUsername,omitempty"`
	ConfigurationPassword string `json:"configurationPassword,omitempty"`
	DockerFile            string `json:"dockerFile,omitempty"`
	ConfigurationFile     string `json:"configurationFile,omitempty"`
	ReportFile            string `json:"reportFile,omitempty"`
}

// HadolintExecuteCommand Executes the Haskell Dockerfile Linter which is a smarter Dockerfile linter that helps you build [best practice](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/) Docker images.
func HadolintExecuteCommand() *cobra.Command {
	const STEP_NAME = "hadolintExecute"

	metadata := hadolintExecuteMetadata()
	var stepConfig hadolintExecuteOptions
	var startTime time.Time

	var createHadolintExecuteCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Executes the Haskell Dockerfile Linter which is a smarter Dockerfile linter that helps you build [best practice](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/) Docker images.",
		Long: `Executes the Haskell Dockerfile Linter which is a smarter Dockerfile linter that helps you build [best practice](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/) Docker images.
The linter is parsing the Dockerfile into an abstract syntax tree (AST) and performs rules on top of the AST.`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.ConfigurationUsername)
			log.RegisterSecret(stepConfig.ConfigurationPassword)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			telemetryData := telemetry.CustomData{}
			telemetryData.ErrorCode = "1"
			handler := func() {
				config.RemoveVaultSecretFiles()
				telemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				telemetryData.ErrorCategory = log.GetErrorCategory().String()
				telemetry.Send(&telemetryData)
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetry.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)
			hadolintExecute(stepConfig, &telemetryData)
			telemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addHadolintExecuteFlags(createHadolintExecuteCmd, &stepConfig)
	return createHadolintExecuteCmd
}

func addHadolintExecuteFlags(cmd *cobra.Command, stepConfig *hadolintExecuteOptions) {
	cmd.Flags().StringVar(&stepConfig.ConfigurationURL, "configurationUrl", os.Getenv("PIPER_configurationUrl"), "URL pointing to the .hadolint.yaml exclude configuration to be used for linting. Also have a look at `configurationFile` which could avoid central configuration download in case the file is part of your repository.")
	cmd.Flags().StringVar(&stepConfig.ConfigurationUsername, "configurationUsername", os.Getenv("PIPER_configurationUsername"), "The username to authenticate")
	cmd.Flags().StringVar(&stepConfig.ConfigurationPassword, "configurationPassword", os.Getenv("PIPER_configurationPassword"), "The password to authenticate")
	cmd.Flags().StringVar(&stepConfig.DockerFile, "dockerFile", `./Dockerfile`, "Dockerfile to be used for the assessment.")
	cmd.Flags().StringVar(&stepConfig.ConfigurationFile, "configurationFile", `.hadolint.yaml`, "Name of the configuration file used locally within the step. If a file with this name is detected as part of your repo downloading the central configuration via `configurationUrl` will be skipped. If you change the file's name make sure your stashing configuration also reflects this.")
	cmd.Flags().StringVar(&stepConfig.ReportFile, "reportFile", `hadolint.xml`, "Name of the result file used locally within the step.")

}

// retrieve step metadata
func hadolintExecuteMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "hadolintExecute",
			Aliases:     []config.Alias{},
			Description: "Executes the Haskell Dockerfile Linter which is a smarter Dockerfile linter that helps you build [best practice](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/) Docker images.",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Parameters: []config.StepParameters{
					{
						Name:        "configurationUrl",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name: "configurationUsername",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "configurationCredentialsId",
								Param: "username",
								Type:  "secret",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
					{
						Name: "configurationPassword",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "configurationCredentialsId",
								Param: "password",
								Type:  "secret",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
					{
						Name:        "dockerFile",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "dockerfile"}},
					},
					{
						Name:        "configurationFile",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "reportFile",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
				},
			},
			Containers: []config.Container{
				{Name: "hadolint", Image: "hadolint/hadolint:latest-debian"},
			},
		},
	}
	return theMetaData
}
