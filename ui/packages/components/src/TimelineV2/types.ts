import type { UserlandSpanType } from '../RunDetailsV3/types';

export type Trace = {
  attempts: number | null;
  childrenSpans?: Trace[];
  endedAt: string | null;
  isRoot: boolean;
  name: string;
  outputID: string | null;
  queuedAt: string;
  spanID: string;
  stepID?: string | null;
  startedAt: string | null;
  status: string;
  stepInfo: StepInfoInvoke | StepInfoSleep | StepInfoWait | StepInfoRun | StepInfoSignal | null;
  stepOp?: string | null;
  isUserland: boolean;
  userlandSpan: UserlandSpanType | null;
};

type StepInfoInvoke = {
  triggeringEventID: string;
  functionID: string;
  timeout: string;
  returnEventID: string | null;
  runID: string | null;
  timedOut: boolean | null;
};

type StepInfoSleep = {
  sleepUntil: string;
};

type StepInfoWait = {
  eventName: string;
  expression: string | null;
  timeout: string;
  foundEventID: string | null;
  timedOut: boolean | null;
};

type StepInfoRun = {
  type: string | null;
};

type StepInfoSignal = {
  signal: string;
  timeout: string;
  timedOut: boolean | null;
};

export function isStepInfoRun(stepInfo: Trace['stepInfo']): stepInfo is StepInfoRun {
  if (!stepInfo) {
    return false;
  }

  return 'type' in stepInfo;
}

export function isStepInfoInvoke(stepInfo: Trace['stepInfo']): stepInfo is StepInfoInvoke {
  if (!stepInfo) {
    return false;
  }

  return 'triggeringEventID' in stepInfo;
}

export function isStepInfoSleep(stepInfo: Trace['stepInfo']): stepInfo is StepInfoSleep {
  if (!stepInfo) {
    return false;
  }

  return 'sleepUntil' in stepInfo;
}

export function isStepInfoWait(stepInfo: Trace['stepInfo']): stepInfo is StepInfoWait {
  if (!stepInfo) {
    return false;
  }

  return 'foundEventID' in stepInfo;
}

export function isStepInfoSignal(stepInfo: Trace['stepInfo']): stepInfo is StepInfoSignal {
  if (!stepInfo) {
    return false;
  }

  return 'signal' in stepInfo;
}
