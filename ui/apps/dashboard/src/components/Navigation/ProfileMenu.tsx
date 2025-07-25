'use client';

import dynamic from 'next/dynamic';
import NextLink from 'next/link';
import { useClerk } from '@clerk/nextjs';
import { Listbox } from '@headlessui/react';
import {
  RiArrowLeftRightLine,
  RiBillLine,
  RiEqualizerLine,
  RiGroupLine,
  RiLogoutCircleLine,
  RiUserLine,
  RiUserSharedLine,
} from '@remixicon/react';

import { pathCreator } from '@/utils/urls';

const ModeSwitch = dynamic(() => import('@inngest/components/ThemeMode/ModeSwitch'), {
  ssr: false,
});

type Props = React.PropsWithChildren<{
  isMarketplace: boolean;
}>;

export const ProfileMenu = ({ children, isMarketplace }: Props) => {
  return (
    <Listbox>
      <Listbox.Button className="w-full cursor-pointer ring-0">{children}</Listbox.Button>
      <div className="relative">
        <Listbox.Options className="bg-canvasBase border-muted shadow-primary absolute -right-48 bottom-4 z-50 ml-8 w-[199px] rounded border ring-0 focus:outline-none">
          <Listbox.Option
            disabled
            value="themeMode"
            className="text-muted mx-2 my-2 flex h-8 items-center justify-between px-2 text-[13px]"
          >
            <div>Theme</div>
            <ModeSwitch />
          </Listbox.Option>

          <hr className="border-subtle" />

          <NextLink href="/settings/user" scroll={false}>
            <Listbox.Option
              className="text-muted hover:bg-canvasSubtle mx-2 mt-2 flex h-8 cursor-pointer items-center px-2 text-[13px]"
              value="userProfile"
            >
              <div className="hover:bg-canvasSubtle flex flex-row items-center justify-start">
                <RiUserLine className="text-muted mr-2 h-4 w-4" />
                <div>Your Profile</div>
              </div>
            </Listbox.Option>
          </NextLink>
          <NextLink href="/settings/organization" scroll={false}>
            <Listbox.Option
              className="text-muted hover:bg-canvasSubtle mx-2 mt-2 flex h-8 cursor-pointer items-center px-2 text-[13px]"
              value="org"
            >
              <div className="hover:bg-canvasSubtle flex flex-row items-center justify-start">
                <RiEqualizerLine className="text-muted mr-2 h-4 w-4 " />
                <div>Your Organization</div>
              </div>
            </Listbox.Option>
          </NextLink>
          <NextLink href="/settings/organization/organization-members" scroll={false}>
            <Listbox.Option
              className="text-muted hover:bg-canvasSubtle mx-2 mt-2 flex h-8 cursor-pointer items-center px-2 text-[13px]"
              value="members"
            >
              <div className="hover:bg-canvasSubtle flex flex-row items-center justify-start">
                <RiGroupLine className="text-muted mr-2 h-4 w-4" />
                <div>Members</div>
              </div>
            </Listbox.Option>
          </NextLink>

          {!isMarketplace && (
            <NextLink href={pathCreator.billing()} scroll={false}>
              <Listbox.Option
                className="text-muted hover:bg-canvasSubtle mx-2 mt-2 flex h-8 cursor-pointer items-center px-2 text-[13px]"
                value="billing"
              >
                <div className="hover:bg-canvasSubtle flex flex-row items-center justify-start">
                  <RiBillLine className="text-muted mr-2 h-4 w-4" />
                  <div>Billing</div>
                </div>
              </Listbox.Option>
            </NextLink>
          )}
          <a href="/organization-list">
            <Listbox.Option
              className="text-muted hover:bg-canvasSubtle m-2 flex h-8 cursor-pointer items-center px-2 text-[13px]"
              value="switchOrg"
            >
              <div className="hover:bg-canvasSubtle flex flex-row items-center justify-start">
                <RiArrowLeftRightLine className="text-muted mr-2 h-4 w-4" />
                <div>Switch Organization</div>
              </div>
            </Listbox.Option>
          </a>

          <hr className="border-subtle" />

          <NextLink href="/sign-in/choose" scroll={false}>
            <Listbox.Option
              className="text-muted hover:bg-canvasSubtle m-2 mx-2 flex h-8 cursor-pointer items-center px-2 text-[13px]"
              value="switchAccount"
            >
              <div className="hover:bg-canvasSubtle flex flex-row items-center justify-start">
                <RiUserSharedLine className="text-muted mr-2 h-4 w-4" />
                <div>Switch Account</div>
              </div>
            </Listbox.Option>
          </NextLink>
          <hr className="border-subtle" />
          <Listbox.Option
            className="text-muted hover:bg-canvasSubtle m-2 flex h-8 cursor-pointer items-center px-2 text-[13px]"
            value="signOut"
          >
            <SignOut isMarketplace={isMarketplace} />
          </Listbox.Option>
        </Listbox.Options>
      </div>
    </Listbox>
  );
};

function SignOut({ isMarketplace }: { isMarketplace: boolean }) {
  const { signOut, session } = useClerk();

  const content = (
    <div className="hover:bg-canvasSubtle flex flex-row items-center justify-start">
      <RiLogoutCircleLine className="text-muted mr-2 h-4 w-4" />
      <div>Sign Out </div>
    </div>
  );

  if (!isMarketplace) {
    // Sign out via Clerk.
    return (
      <button
        onClick={async () => {
          await signOut({ sessionId: session?.id, redirectUrl: '/sign-in/choose' });
        }}
      >
        {content}
      </button>
    );
  }

  // Sign out via our backend.
  return <NextLink href={`${process.env.NEXT_PUBLIC_API_URL}/v1/logout`}>{content}</NextLink>;
}
