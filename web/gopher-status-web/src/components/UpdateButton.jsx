import { ActionIcon } from '@mantine/core';
import { IconRefresh } from '@tabler/icons-react';

export function UpdateButton({ onUpdate, loading }) {
    return (
        <ActionIcon
            onClick={onUpdate}
            variant="default"
            size="lg"
            radius="xl"
            aria-label="Update list"
            loading={loading}
        >
            <IconRefresh size={20} />
        </ActionIcon>
    );
}
