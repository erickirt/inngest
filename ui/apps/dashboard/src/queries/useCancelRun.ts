import { useMutation } from 'urql';

import { graphql } from '@/gql';

const mutation = graphql(`
  mutation CancelRun($envID: UUID!, $runID: ULID!) {
    cancelRun(envID: $envID, runID: $runID) {
      id
    }
  }
`);

export function useCancelRun({ envID, runID }: { envID: string; runID: string }) {
  const [, mutate] = useMutation(mutation);

  return async () => {
    const res = await mutate({ envID, runID });
    if (res.error) {
      // Throw error so that the modal can catch and display it
      throw res.error;
    }
  };
}