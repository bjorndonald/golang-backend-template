import React from "react";
import { Metadata } from "next";
import AuthLayout from "@/components/Layouts/AuthLayout";
import ResetPasswordForm from "@/app/auth/reset-password/ResetPasswordForm";

export const metadata: Metadata = {
    title: "Next.js SignIn Page | TailAdmin - Next.js Dashboard Template",
    description: "This is Next.js Signin Page TailAdmin Dashboard Template",
};

interface ResetPasswordProps {
    searchParams?: Promise<{
        email?: string;
    }>
}

const ResetPassword: React.FC<ResetPasswordProps> = async ({searchParams}) => {
    const email = await(await searchParams)?.email
    return (
        <AuthLayout>
            {/* <Breadcrumb pageName="Sign In" /> */}
            <div className="my-8">
                <div className="rounded-sm border max-w-xl mx-auto px-4 border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
                    <div className="flex flex-wrap items-center">
                        <div className="w-full ">
                            <div className="w-full p-4 sm:p-12.5 xl:p-17.5">

                                <h2 className="mb-9 text-2xl font-bold text-black dark:text-white sm:text-title-xl2">
                                    Reset Password
                                </h2>

                                <ResetPasswordForm email={email} />
                            </div>
                        </div>
                    </div>
                </div>

            </div>
        </AuthLayout>
    );
};

export default ResetPassword;
