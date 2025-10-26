CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    read_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE notifications IS 'Таблица для хранения уведомлений пользователей';
COMMENT ON COLUMN notifications.user_id IS 'ID пользователя, которому отправлено уведомление';
COMMENT ON COLUMN notifications.type IS 'Тип уведомления (registration_email, chat_message, etc.)';
COMMENT ON COLUMN notifications.status IS 'Статус уведомления (pending, sent, failed)';
COMMENT ON COLUMN notifications.metadata IS 'Дополнительные данные в формате JSON';
