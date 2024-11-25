"use client";
import apiInstance from '@/services/api';
import { zodResolver } from '@hookform/resolvers/zod';
import React, { useEffect } from 'react'
import { useForm, FieldValues } from 'react-hook-form';
import toast from 'react-hot-toast';
import { z } from 'zod';
import { AxiosError } from 'axios';
import Cookies from 'js-cookie'

const schema = z.object({
    otp: z.string()
})

type TwoFAInput = z.infer<typeof schema>

const TwoFAForm = ({ token }: { token: string | undefined }) => {
    const methods = useForm<TwoFAInput>({
        resolver: zodResolver(schema),
    });
    const {
        setValue,
        handleSubmit,
        register,
    } = methods;

    useEffect(() => {
        if (!token) window.location.replace("/auth/signin")
       
        return () => {

        }
    }, [])


    const onSubmit = async (data: FieldValues) => {
        try {
            toast.loading("Loading...", { id: "loading" })
            const response = await apiInstance.post('/api/v1/auth/2fa/verify/'+ token, data);

            toast.success(response.data.message)
            Cookies.set('accessToken', response.data.data.access_token);
            toast.remove("loading")
            setTimeout(() => {
                window.location.replace(`/`)
            }, 2000);
        } catch (error) {
            toast.remove("loading")
            if (error instanceof AxiosError) {
                toast.error(error.response?.data.message)
            }
        }
    }

    const resendEmail = async () => {
        try {
            toast.loading("Loading...", { id: "loading" })
            const response = await apiInstance.post('/api/v1/auth/2fa/'+token);
            toast.success(response.data.message)
            toast.remove("loading")
            setTimeout(() => {
                window.location.replace("/auth/2fa?token=" + token)
            })
        } catch (error) {
            toast.remove("loading")
            if (error instanceof AxiosError) {
                if (error.response?.data?.error.includes("token is expired")){
                    toast.error("Token has expired. Please try logging in again.")
                    setTimeout(() => {
                        window.location.replace(`/auth/signin`)
                    }, 3000);
                    return
                }
                    
                toast.error(error.response?.data?.message)

                
            }

        }

    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            

            <div className="mb-6">
                <label className="mb-2.5 block font-medium text-black dark:text-white">
                    OTP
                </label>
                <div className="flex outline-none border-form-strokedark focus:border-primary focus-visible:shadow-none dark:border-form-strokedark dark:bg-form-input dark:text-white dark:focus:border-primary join items-center">
                    <input
                        {...register("otp")}
                        type="text"
                        placeholder="6+ Characters, 1 Capital letter"
                        className="w-full rounded-lg outline-none border-none bg-transparent py-4 pl-6 pr-10 text-black"
                    />

                    <a onClick={resendEmail} className=" cursor-pointer rounded-lg border border-primary text-xs bg-primary p-4 text-white transition hover:bg-opacity-90">
                        Resend
                    </a>
                </div>
            </div>

            <div className="mb-5">
                <button
                    type="submit"
                    className="w-full cursor-pointer rounded-lg border border-primary bg-primary p-4 text-white transition hover:bg-opacity-90"
                >
                    Confirm
                </button>
            </div>

            
        </form>
    )
}

export default TwoFAForm