"use client";
import PasswordIcon from '@/components/icons/Password.icon';
import apiInstance from '@/services/api';
import { zodResolver } from '@hookform/resolvers/zod';
import { AxiosError } from 'axios';
import React from 'react'
import { useForm, FieldValues } from 'react-hook-form';
import toast from 'react-hot-toast';
import { z } from 'zod';

const schema = z.object({
    password: z.string(),
    passwordConfirm: z.string()
}).refine((x) => x.password == x.passwordConfirm, "The passwords do not match")

type ChangePasswordInput = z.infer<typeof schema>

const ChangePasswordForm = ({ reset_token }: { reset_token:string|undefined}) => {
    
    const methods = useForm<ChangePasswordInput>({
        resolver: zodResolver(schema),
    });
    const {
        handleSubmit,
        formState: {errors},
        register,
    } = methods;

    const onSubmit = async (data: FieldValues) => {
        try {
            toast.loading("Loading...", { id: "loading" })
            const response = await apiInstance.post(`/api/v1/auth/reset-password/confirm/${reset_token}`, data);
            toast.success(response.data.message)
            toast.remove("loading")
            window.location.replace("/auth/signin")
        } catch (error) {
            toast.remove("loading")
            if (error instanceof AxiosError) {
                if (!!error.response && error.response?.status >= 400)
                    toast.error(error.response?.data?.message)
            } else toast.error("Server error")
        }
        
    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>

            <div className="mb-4">
                <label className="mb-2.5 block font-medium text-black dark:text-white">
                    Password
                </label>
                <div className="relative">
                    <input
                        {...register("password")}
                        type="password"
                        placeholder="Enter your Password"
                        className="w-full rounded-lg border border-stroke bg-transparent py-4 pl-6 pr-10 text-black outline-none focus:border-primary focus-visible:shadow-none dark:border-form-strokedark dark:bg-form-input dark:text-white dark:focus:border-primary"
                    />

                    <span className="absolute right-4 top-4">
                        <PasswordIcon />
                    </span>
                </div>
                {!!errors.password && <p className='text-red-500'>{errors.password?.message}</p>}
                
            </div>

            <div className="mb-6">
                <label className="mb-2.5 block font-medium text-black dark:text-white">
                    Confirm Password
                </label>
                <div className="relative">
                    <input
                        {...register("passwordConfirm")}
                        type="password"
                        placeholder="6+ Characters, 1 Capital letter"
                        className="w-full rounded-lg border border-stroke bg-transparent py-4 pl-6 pr-10 text-black outline-none focus:border-primary focus-visible:shadow-none dark:border-form-strokedark dark:bg-form-input dark:text-white dark:focus:border-primary"
                    />

                    <span className="absolute right-4 top-4">
                        <PasswordIcon />
                    </span>
                </div>
                {!!errors.passwordConfirm && <p className='text-red-500'>{errors.passwordConfirm?.message}</p>}
            </div>

            <div className="mb-5">
                <button
                    type="submit"
                    className="w-full cursor-pointer rounded-lg border border-primary bg-primary p-4 text-white transition hover:bg-opacity-90"
                >
                    Change Password
                </button>
            </div>


        </form>
    )
}

export default ChangePasswordForm