package logger

import (
	"errors"
	"github.com/graniticio/granitic/config"
	"github.com/graniticio/granitic/facility/decorator"
	"github.com/graniticio/granitic/instance"
	"github.com/graniticio/granitic/ioc"
	"github.com/graniticio/granitic/logging"
)

const applicationLoggingDecoratorName = instance.FrameworkPrefix + "ApplicationLoggingDecorator"
const applicationLoggingManagerName = instance.FrameworkPrefix + "ApplicationLoggingManager"

type ApplicationLoggingFacilityBuilder struct {
}

func (alfb *ApplicationLoggingFacilityBuilder) BuildAndRegister(lm *logging.ComponentLoggerManager, ca *config.ConfigAccessor, cn *ioc.ComponentContainer) error {
	defaultLogLevelLabel, err := ca.StringVal("ApplicationLogger.DefaultLogLevel")

	if err != nil {
		return alfb.error(err.Error())
	}

	defaultLogLevel, err := logging.LogLevelFromLabel(defaultLogLevelLabel)

	if err != nil {
		return alfb.error(err.Error())
	}

	initialLogLevelsByComponent := ca.ObjectVal("ApplicationLogger.ComponentLogLevels")

	writers, err := alfb.buildWriters(ca)
	formatter, err := alfb.buildFormatter(ca)

	if err != nil {
		return alfb.error(err.Error())
	}

	//Update the bootstrapped framework logger with the newly configured writers and formatter
	lm.UpdateWritersAndFormatter(writers, formatter)

	alm := logging.CreateComponentLoggerManager(defaultLogLevel, initialLogLevelsByComponent, writers, formatter)
	cn.WrapAndAddProto(applicationLoggingManagerName, alm)

	ald := new(decorator.ApplicationLogDecorator)
	ald.LoggerManager = alm
	ald.FrameworkLogger = lm.CreateLogger(applicationLoggingDecoratorName)

	cn.WrapAndAddProto(applicationLoggingDecoratorName, ald)

	return nil
}

func (alfb *ApplicationLoggingFacilityBuilder) buildFormatter(ca *config.ConfigAccessor) (*logging.LogMessageFormatter, error) {

	lmf := new(logging.LogMessageFormatter)

	err := ca.Populate("LogWriting.PrefixFormat", lmf)

	if err != nil {
		return nil, err
	}

	if lmf.PrefixFormat == "" && lmf.PrefixPreset == "" {
		lmf.PrefixPreset = logging.FrameworkPresetPrefix
	}

	err = lmf.Init()

	return lmf, err

}

func (alfb *ApplicationLoggingFacilityBuilder) buildWriters(ca *config.ConfigAccessor) ([]logging.LogWriter, error) {
	writers := make([]logging.LogWriter, 0)

	console, err := ca.BoolVal("LogWriting.EnableConsoleLogging")

	if err != nil {
		return nil, err
	}

	if console {
		writers = append(writers, new(logging.ConsoleWriter))
	}

	file, err := ca.BoolVal("LogWriting.EnableFileLogging")

	if err != nil {
		return nil, err
	}

	if file {
		fileWriter := new(logging.AsynchFileWriter)

		err := ca.Populate("LogWriting.File", fileWriter)

		if err != nil {
			return nil, err
		}

		err = fileWriter.Init()

		if err != nil {
			return nil, err
		}

		writers = append(writers, fileWriter)
	}

	return writers, nil
}

func (alfb *ApplicationLoggingFacilityBuilder) error(suffix string) error {

	return errors.New("Unable to initialise application logging: " + suffix)

}

func (alfb *ApplicationLoggingFacilityBuilder) FacilityName() string {
	return "ApplicationLogging"
}

func (alfb *ApplicationLoggingFacilityBuilder) DependsOnFacilities() []string {
	return []string{}
}
