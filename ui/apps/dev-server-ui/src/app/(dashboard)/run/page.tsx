'use client';

import { RunDetailsV2 } from '@inngest/components/RunDetailsV2/RunDetailsV2';
import { RunDetailsV3 } from '@inngest/components/RunDetailsV3/RunDetailsV3';
import { useLegacyTrace } from '@inngest/components/SharedContext/useLegacyTrace';
import { useSearchParam } from '@inngest/components/hooks/useSearchParam';
import { cn } from '@inngest/components/utils/classNames';

import { useGetRun } from '@/hooks/useGetRun';
import { useGetTraceResult } from '@/hooks/useGetTraceResult';
import { useGetTrigger } from '@/hooks/useGetTrigger';
import { pathCreator } from '@/utils/pathCreator';

export default function Page() {
  const [runID] = useSearchParam('runID');
  const getRun = useGetRun();
  const getTraceResult = useGetTraceResult();
  const getTrigger = useGetTrigger();

  const traceAIEnabled = true;
  const { enabled: legacyTraceEnabled } = useLegacyTrace();

  if (!runID) {
    throw new Error('missing runID in search params');
  }

  return (
    <div className={cn('bg-canvasBase overflow-y-auto pt-8')}>
      {traceAIEnabled && !legacyTraceEnabled ? (
        <RunDetailsV3
          pathCreator={pathCreator}
          standalone
          getResult={getTraceResult}
          getRun={getRun}
          getTrigger={getTrigger}
          pollInterval={2500}
          runID={runID}
        />
      ) : (
        <RunDetailsV2
          pathCreator={pathCreator}
          standalone
          getResult={getTraceResult}
          getRun={getRun}
          getTrigger={getTrigger}
          pollInterval={2500}
          runID={runID}
          traceAIEnabled={traceAIEnabled}
        />
      )}
    </div>
  );
}
