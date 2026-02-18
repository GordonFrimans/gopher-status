// ThemeToggle.jsx
import { ActionIcon, useMantineColorScheme } from '@mantine/core';
import { IconSun, IconMoon } from '@tabler/icons-react';

export function ThemeToggle() {
    const { colorScheme, toggleColorScheme } = useMantineColorScheme();

    return (
        <ActionIcon
            onClick={toggleColorScheme}
            variant="default"
            size="lg"
            radius="xl"
            aria-label="Toggle color scheme"


        >
            {colorScheme === 'dark' ? <IconSun size={20} /> : <IconMoon size={20} />}
        </ActionIcon>
    );
}
