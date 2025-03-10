'use client';

import { useCallback, useEffect, useState } from 'react';
import DescriptionListItem from '@inngest/components/Apps/DescriptionListItem';
import { Button } from '@inngest/components/Button';
import { Pill } from '@inngest/components/Pill/Pill';
import {
  ConnectV1WorkerConnectionsOrderByDirection,
  ConnectV1WorkerConnectionsOrderByField,
  type ConnectV1WorkerConnectionsOrderBy,
  type GroupedWorkerStatus,
  type PageInfo,
  type Worker,
  type WorkerStatus,
} from '@inngest/components/types/workers';
import { transformLanguage } from '@inngest/components/utils/appsParser';
import { convertGroupedWorkerStatusToWorkerStatuses } from '@inngest/components/utils/workerParser';
import { RiArrowLeftSLine, RiArrowRightSLine } from '@remixicon/react';
import { keepPreviousData, useQuery } from '@tanstack/react-query';
import { type Row, type SortingState } from '@tanstack/react-table';

import CompactPaginatedTable from '../Table/CompactPaginatedTable';
import WorkerStatusFilter from './WorkerStatusFilter';
import { useColumns } from './columns';

const columnToTimeField: Record<string, ConnectV1WorkerConnectionsOrderByField> = {
  connectedAt: ConnectV1WorkerConnectionsOrderByField.ConnectedAt,
  disconnectedAt: ConnectV1WorkerConnectionsOrderByField.DisconnectedAt,
  lastHeartbeatAt: ConnectV1WorkerConnectionsOrderByField.LastHeartbeatAt,
};

const refreshInterval = 5000;

export function WorkersTable({
  appID,
  getWorkers,
  getWorkerCount,
}: {
  appID: string;
  getWorkerCount: ({ appID }: { appID: string }) => Promise<number>;
  getWorkers: ({
    appID,
    orderBy,
    cursor,
    pageSize,
    status,
  }: {
    appID: string;
    orderBy: ConnectV1WorkerConnectionsOrderBy[];
    cursor: string | null;
    pageSize: number;
    status: WorkerStatus[];
  }) => Promise<{ workers: Worker[]; pageInfo: PageInfo; totalCount: number }>;
}) {
  const columns = useColumns();
  const [sorting, setSorting] = useState<SortingState>([
    {
      id: 'connectedAt',
      desc: true,
    },
  ]);

  const [orderBy, setOrderBy] = useState<ConnectV1WorkerConnectionsOrderBy[]>([
    {
      field: ConnectV1WorkerConnectionsOrderByField.ConnectedAt,
      direction: ConnectV1WorkerConnectionsOrderByDirection.Asc,
    },
  ]);
  const [filteredStatus, setFilteredStatus] = useState<WorkerStatus[]>([]);
  const [cursor, setCursor] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);

  const onStatusFilterChange = useCallback(
    (value: GroupedWorkerStatus[]) => {
      const workerStatuses = value.flatMap(convertGroupedWorkerStatusToWorkerStatuses);
      setFilteredStatus(workerStatuses);
      // Back to first page when we sort changes
      setCursor(null);
      setPage(1);
    },
    [setFilteredStatus]
  );

  const {
    isPending, // first load, no data
    error,
    data: workerConnsData,
    isFetching, // refetching
  } = useQuery({
    queryKey: ['worker-connections', { appID, orderBy, cursor, pageSize, status: filteredStatus }],
    queryFn: useCallback(() => {
      return getWorkers({ appID, orderBy, cursor, pageSize, status: filteredStatus });
    }, [getWorkers, appID, orderBy, cursor, pageSize, filteredStatus]),
    placeholderData: keepPreviousData,
    refetchInterval: !cursor || page === 1 ? refreshInterval : 0,
  });
  const pageInfo = workerConnsData?.pageInfo;

  const { data: totalCount } = useQuery({
    queryKey: ['worker-count', { appID }],
    queryFn: useCallback(() => {
      return getWorkerCount({ appID });
    }, [getWorkerCount, appID]),
    placeholderData: keepPreviousData,
  });

  const numberOfPages = Math.ceil((totalCount || 0) / pageSize);

  useEffect(() => {
    const sortEntry = sorting[0];
    if (!sortEntry) return;

    const sortColumn = sortEntry.id;
    if (sortColumn && columnToTimeField[sortColumn]) {
      const orderBy: ConnectV1WorkerConnectionsOrderBy[] = [
        {
          field:
            columnToTimeField[sortColumn] ?? ConnectV1WorkerConnectionsOrderByField.ConnectedAt,
          direction: sortEntry.desc
            ? ConnectV1WorkerConnectionsOrderByDirection.Desc
            : ConnectV1WorkerConnectionsOrderByDirection.Asc,
        },
      ];
      setOrderBy(orderBy);
      // Back to first page when we sort changes
      setCursor(null);
      setPage(1);
    }
  }, [sorting, setOrderBy]);

  return (
    <div>
      <h4 className="text-subtle mb-4 text-xl">Workers</h4>
      <div className="mb-4 flex items-center">
        <WorkerStatusFilter
          selectedStatuses={filteredStatus}
          onStatusesChange={onStatusFilterChange}
        />
      </div>
      <CompactPaginatedTable
        columns={columns}
        data={workerConnsData?.workers || []}
        isLoading={isPending}
        sorting={sorting}
        setSorting={setSorting}
        enableExpanding={true}
        renderSubComponent={SubComponent}
        getRowCanExpand={() => true}
        footer={
          (workerConnsData?.totalCount ?? 0) > pageSize ? (
            <div className="flex items-center justify-end gap-2 px-6 py-3">
              <Button
                kind="secondary"
                appearance="outlined"
                disabled={page === 1}
                // disabled={!pageInfo?.hasPreviousPage} TODO: use this once it's fixed in the BE
                icon={<RiArrowLeftSLine />}
                onClick={() => {
                  setCursor(pageInfo?.startCursor || null);
                  setPage(page - 1);
                }}
              />
              {page}/{numberOfPages}
              <Button
                kind="secondary"
                appearance="outlined"
                disabled={!pageInfo?.hasNextPage}
                icon={<RiArrowRightSLine />}
                onClick={() => {
                  setCursor(pageInfo?.endCursor || null);
                  setPage(page + 1);
                }}
              />
            </div>
          ) : undefined
        }
      />
    </div>
  );
}

function SubComponent({ row }: { row: Row<Worker> }) {
  return (
    <dl className="bg-canvasSubtle mx-9 mb-6 mt-[10px] grid grid-cols-4 gap-2 p-4">
      <DescriptionListItem term="SDK version" detail={row.original.sdkVersion} />
      <DescriptionListItem term="SDK language" detail={transformLanguage(row.original.sdkLang)} />
      <DescriptionListItem term="No. of functions" detail={row.original.functionCount.toString()} />
      <DescriptionListItem
        term="System attributes"
        detail={
          <div className="flex items-center gap-1">
            <Pill>{row.original.os + ' OS'}</Pill>
            <Pill>{row.original.cpuCores + ' CPU cores'}</Pill>
          </div>
        }
      />
      <DescriptionListItem className="col-span-3" term="Worker IP" detail={row.original.workerIp} />
    </dl>
  );
}
