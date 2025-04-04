import type { Route } from 'next';
import AppDetailsCard from '@inngest/components/Apps/AppDetailsCard';
import { Link } from '@inngest/components/Link';

type Props = {
  sync: {
    platform?: string | null;
    vercelDeploymentID: string | null;
    vercelDeploymentURL: string | null;
    vercelProjectID: string | null;
    vercelProjectURL: string | null;
  };
};

export function PlatformSection({ sync }: Props) {
  const { platform, vercelDeploymentID, vercelDeploymentURL, vercelProjectID, vercelProjectURL } =
    sync;
  if (!platform) {
    return null;
  }

  let deploymentValue;
  if (vercelDeploymentID && vercelDeploymentURL) {
    deploymentValue = (
      <Link href={vercelDeploymentURL as Route} target="_blank" size="small">
        <span className="truncate">{vercelDeploymentID}</span>
      </Link>
    );
  } else {
    deploymentValue = '-';
  }

  let projectValue;
  if (vercelProjectID && vercelProjectURL) {
    projectValue = (
      <Link href={vercelProjectURL as Route} target="_blank" size="small">
        <span className="truncate">{vercelProjectID}</span>
      </Link>
    );
  } else {
    projectValue = '-';
  }

  return (
    <>
      <AppDetailsCard.Item detail={<div className="truncate">{platform}</div>} term="Platform" />
      <AppDetailsCard.Item detail={projectValue} term="Vercel project" />
      <AppDetailsCard.Item detail={deploymentValue} term="Vercel deployment" />
    </>
  );
}
