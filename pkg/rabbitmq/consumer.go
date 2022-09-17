package rabbitmq

import (
	"github.com/streadway/amqp"

	"github.com/art-injener/iot-platform/internal/config"
	"github.com/art-injener/iot-platform/pkg/logger"
)

type Consumer struct {
	Config         *config.RabbitConfig
	MessageChannel <-chan amqp.Delivery
	Logger         *logger.Logger
	connection     *amqp.Connection
	amqpChannel    *amqp.Channel
}

func (c *Consumer) Initialize() error {
	var err error
	c.connection, err = amqp.Dial(c.Config.Url)
	if err != nil {
		c.Logger.Error().Msgf("Can't connect to AMQP by url %s", c.Config.Url)
		return err
	}

	c.amqpChannel, err = c.connection.Channel()
	if err != nil {
		c.Logger.Error().Msg("Can't create a amqpChannel")
		return err
	}

	queue, err := c.amqpChannel.QueueDeclare("device_info", true, false, false, false, nil)
	if err != nil {
		c.Logger.Error().Msgf("Could not declare `add` queue %s", c.Config.Queue.QueueName)
		return err
	}

	err = c.amqpChannel.Qos(
		c.Config.Qos.PrefetchCount,
		c.Config.Qos.PrefetchSize,
		c.Config.Qos.Global,
	)
	if err != nil {
		c.Logger.Error().Msg("Could not configure QoS")
		return err
	}

	c.MessageChannel, err = c.amqpChannel.Consume(
		queue.Name,
		c.Config.Consumer.Tag,
		c.Config.Consumer.AutoAck,
		c.Config.Consumer.Exclusive,
		c.Config.Consumer.NoLocal,
		c.Config.Consumer.NoWait,
		nil,
	)
	if err != nil {
		c.Logger.Error().Msgf("Could not register consumer with queue name %s", queue.Name)
		return err
	}

	return nil
}

func (c *Consumer) Stop() {
	c.amqpChannel.Close()
	c.connection.Close()
}
