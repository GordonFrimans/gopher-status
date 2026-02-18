import { useMemo } from 'react';
import { Table, Badge, ActionIcon, Group, Text } from '@mantine/core';
import { IconTrash, IconExternalLink } from '@tabler/icons-react';

const getStatusBadge = (status) => {
    switch (status) {
        case 'UP':
            return <Badge color="green" variant="light">UP</Badge>;
        case 'DOWN':
            return <Badge color="red" variant="filled">DOWN</Badge>;
        case 'PENDING':
            return <Badge color="gray" variant="outline">PENDING</Badge>;
        default:
            return <Badge color="gray">{status}</Badge>;
    }
};

export function MonitorTable({ monitors = [], onDelete }) {
    const sortedMonitors = useMemo(
        () => [...monitors].sort((a, b) => Number(a.id) - Number(b.id)),
        [monitors]
    );

    const rows = sortedMonitors.map((monitor) => (
        <Table.Tr key={monitor.id}>
            <Table.Td>
                <Text size="sm" c="dimmed">#{monitor.id}</Text>
            </Table.Td>

            <Table.Td style={{ fontWeight: 500 }}>
                {monitor.name}
            </Table.Td>

            <Table.Td>
                <Group gap="xs" wrap="nowrap">
                    <Text
                        size="sm"
                        c="blue"
                        component="a"
                        href={monitor.url}
                        target="_blank"
                        style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}
                    >
                        {monitor.url}
                    </Text>
                    <ActionIcon size="xs" variant="subtle" color="gray" component="a" href={monitor.url} target="_blank">
                        <IconExternalLink size={12} />
                    </ActionIcon>
                </Group>
            </Table.Td>

            <Table.Td>{getStatusBadge(monitor.status)}</Table.Td>

            <Table.Td>
                <Text size="sm">{monitor.interval}s</Text>
            </Table.Td>

            <Table.Td>
                <Text size="sm" c="dimmed" ff="monospace">
                    {monitor.lastCheck ? monitor.lastCheck : '-'}
                </Text>
            </Table.Td>

            <Table.Td style={{ textAlign: 'center' }}>
                <ActionIcon
                    variant="light"
                    color="red"
                    onClick={() => onDelete(monitor.id)}
                    aria-label="Delete monitor"
                >
                    <IconTrash size={16} />
                </ActionIcon>
            </Table.Td>
        </Table.Tr>
    ));

    return (
        <Table striped highlightOnHover withTableBorder style={{ tableLayout: 'fixed' }}>
            <Table.Thead>
                <Table.Tr>
                    <Table.Th style={{ width: 70 }}>ID</Table.Th>
                    <Table.Th>Name</Table.Th>
                    <Table.Th>URL</Table.Th>
                    <Table.Th style={{ width: 110 }}>Status</Table.Th>
                    <Table.Th style={{ width: 90 }}>Interval</Table.Th>
                    <Table.Th style={{ width: 190 }}>Last Check</Table.Th>
                    <Table.Th style={{ width: 90, textAlign: 'center' }}>Actions</Table.Th>
                </Table.Tr>
            </Table.Thead>
            <Table.Tbody>{rows}</Table.Tbody>
        </Table>
    );
}
