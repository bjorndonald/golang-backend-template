"use client";
import EmailIcon from '@/components/icons/Email.icon';
import PasswordIcon from '@/components/icons/Password.icon';
import apiInstance from '@/services/api';
import { zodResolver } from '@hookform/resolvers/zod';
import Link from 'next/link';
import React, { useEffect } from 'react'
import { useForm, FieldValues } from 'react-hook-form';
import toast from 'react-hot-toast';
import { z } from 'zod';
import Cookies from 'js-cookie'
import { AxiosError } from 'axios';

const schema = z.object({
    email: z.string().email(),
    otp: z.string()
})

type ResetPasswordInput = z.infer<typeof schema>

const ResetPasswordForm = ({ email }: { email: string | undefined }) => {
    const methods = useForm<ResetPasswordInput>({
        resolver: zodResolver(schema),
    });
    const {
        setValue,
        handleSubmit,
        register,
    } = methods;

    useEffect(() => {
        if(!email) window.location.replace("/auth/forgot-password")
        setValue("email", email ?? "")
    
      return () => {
        
      }
    }, [])
    

    const onSubmit = async (data: FieldValues) => {
        try {
            toast.loading("Loading...", { id: "loading" })
            const response = await apiInstance.post('/api/v1/auth/forgot-password/verify', data);

            toast.success(response.data.message)
            toast.remove("loading")
            // setTimeout(() => {
            //     window.location.replace(`/auth/change-password?reset_token=${response.data.data.access_token}`)
            // }, 3000);
        } catch (error) {
            toast.remove("loading")
            if (error instanceof AxiosError) {
                if (!!error.response && error.response?.status >= 400)
                    toast.error(error.response?.data?.message)
            } else toast.error("Server error")
        }
    }

    const resendEmail = async () => {
        try {
            toast.loading("Loading...", { id: "loading" })
            const response = await apiInstance.post('/api/v1/auth/forgot-password', {email: email});
            toast.success(response.data.message)
            toast.remove("loading")
            setTimeout(() => {
                window.location.replace("/auth/reset-password?email=" + email)
            })
        } catch (error) {
            toast.remove("loading")
            if (error instanceof AxiosError) {
                toast.error(error.response?.data?.message)
            }

        }

    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mb-4">
                <label className="mb-2.5 block font-medium text-black dark:text-white">
                    Email
                </label>
                <div className="relative">
                    <input
                        {...register("email")}
                        type="email"
                        placeholder="Enter your email"
                        className="w-full rounded-lg border border-stroke bg-transparent py-4 pl-6 pr-10 text-black outline-none focus:border-primary focus-visible:shadow-none dark:border-form-strokedark dark:bg-form-input dark:text-white dark:focus:border-primary"
                    />

                    <span className="absolute right-4 top-4">
                        <EmailIcon />
                    </span>
                </div>
            </div>

            <div className="mb-6">
                <label className="mb-2.5 block font-medium text-black dark:text-white">
                    OTP
                </label>
                <div className="relative">
                    <input
                        {...register("otp")}
                        type="password"
                        placeholder="6+ Characters, 1 Capital letter"
                        className="w-full rounded-lg border border-stroke bg-transparent py-4 pl-6 pr-10 text-black outline-none focus:border-primary focus-visible:shadow-none dark:border-form-strokedark dark:bg-form-input dark:text-white dark:focus:border-primary"
                    />

                    <span className="absolute right-4 top-4">
                        <PasswordIcon />
                    </span>
                </div>
            </div>

            <div className="mb-5">
                <button
                    type="submit"
                    className="w-full cursor-pointer rounded-lg border border-primary bg-primary p-4 text-white transition hover:bg-opacity-90" 
                >
                    Reset Password
                </button>
            </div>

            <a onClick={resendEmail} className='text-primary cursor-pointer text-lg'>Resend email</a>
        </form>
    )
}

export default ResetPasswordForm