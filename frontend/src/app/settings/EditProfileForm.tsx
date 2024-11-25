"use client";
import EmailIcon from '@/components/icons/Email.icon'
import PasswordIcon from '@/components/icons/Password.icon'
import React, { useActionState, useContext, useEffect } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { FieldValues, useForm } from 'react-hook-form';
import Link from 'next/link';
import toast from 'react-hot-toast';
import apiInstance from '@/services/api';
import Cookies from 'js-cookie'
import { AxiosError } from 'axios';
import { Edit, Mail, User } from 'lucide-react';
import { AppContext } from '@/components/Layouts/DefaultLayout';

const schema = z.object({
    first_name: z.string(),
    last_name: z.string(),
    email: z.string().email(),
    phone_number: z.string(),
    bio: z.string()
})

type EditProfileInput = z.infer<typeof schema>

const EditProfileForm = () => {
    const {user} = useContext(AppContext)
    const methods = useForm<EditProfileInput>({
        resolver: zodResolver(schema),
    });
    const {
        handleSubmit,
        register,
        setValue
    } = methods;

    useEffect(() => {
      if(!!user) {
        setValue("first_name", user.first_name)
          setValue("last_name", user.last_name)
          setValue("phone_number", user.phone_number)
          setValue("email", user.email)
          setValue("bio", user.bio)
      }
    
      return () => {
      }
    }, [user])
    

    const onSubmit = async (data: FieldValues) => {
        try {
            toast.loading("Loading...", { id: "loading" })
            const response = await apiInstance.put('/api/v1/user/profile', data);

            toast.success(response.data.message)
            toast.remove("loading")
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
            <div className="mb-5.5 flex flex-col gap-5.5 sm:flex-row">
                <div className="w-full sm:w-1/2">
                    <label
                        className="mb-3 block text-sm font-medium text-black dark:text-white"
                        htmlFor="first_name"
                    >
                        First Name
                    </label>
                    <div className="relative">
                        <span className="absolute left-4.5 top-4">
                            <User />
                        </span>
                        <input
                            className="w-full rounded border border-stroke bg-gray py-3 pl-11.5 pr-4.5 text-black focus:border-primary focus-visible:outline-none dark:border-strokedark dark:bg-meta-4 dark:text-white dark:focus:border-primary"
                            type="text"
                           
                            placeholder="Your First Name"
                            {...register("first_name")}
                        />
                    </div>
                </div>

                <div className="w-full sm:w-1/2">
                    <label
                        className="mb-3 block text-sm font-medium text-black dark:text-white"
                        htmlFor="last_name"
                    >
                        Last Name
                    </label>
                    <input
                        className="w-full rounded border border-stroke bg-gray px-4.5 py-3 text-black focus:border-primary focus-visible:outline-none dark:border-strokedark dark:bg-meta-4 dark:text-white dark:focus:border-primary"
                        type="text"
                        {...register("last_name")}
                        placeholder="Your last name here"
                    />
                </div>
            </div>

            <div className="mb-5.5">
                <label
                    className="mb-3 block text-sm font-medium text-black dark:text-white"
                    htmlFor="email"
                >
                    Email Address
                </label>
                <div className="relative">
                    <span className="absolute left-4.5 top-4">
                        <Mail />
                    </span>
                    <input
                        className="w-full rounded border border-stroke bg-gray py-3 pl-11.5 pr-4.5 text-black focus:border-primary focus-visible:outline-none dark:border-strokedark dark:bg-meta-4 dark:text-white dark:focus:border-primary"
                        type="email"
                        disabled
                        {...register("email")}
                        placeholder="Your email address here"
                    />
                </div>
            </div>

            <div className="mb-5.5">
                <label
                    className="mb-3 block text-sm font-medium text-black dark:text-white"
                    htmlFor="phone_number"
                >
                    Phone Number
                </label>
                <input
                    className="w-full rounded border border-stroke bg-gray px-4.5 py-3 text-black focus:border-primary focus-visible:outline-none dark:border-strokedark dark:bg-meta-4 dark:text-white dark:focus:border-primary"
                    type="text"
                    {...register("phone_number")}
                    placeholder="Your phone number here"
                />
            </div>

            <div className="mb-5.5">
                <label
                    className="mb-3 block text-sm font-medium text-black dark:text-white"
                    htmlFor="bio"
                >
                    Bio
                </label>
                <div className="relative">
                    <span className="absolute left-4.5 top-4">
                        <Edit />
                    </span>

                    <textarea
                        className="w-full rounded border border-stroke bg-gray py-3 pl-11.5 pr-4.5 text-black focus:border-primary focus-visible:outline-none dark:border-strokedark dark:bg-meta-4 dark:text-white dark:focus:border-primary"
                        
                        rows={6}
                        placeholder="Write your bio here"
                        {...register("bio")}
                    ></textarea>
                </div>
            </div>

            <div className="flex justify-end gap-4.5">
                <button
                    className="flex justify-center rounded border border-stroke px-6 py-2 font-medium text-black hover:shadow-1 dark:border-strokedark dark:text-white"
                    type="reset"
                >
                    Cancel
                </button>
                <button
                    className="flex justify-center rounded bg-primary px-6 py-2 font-medium text-gray hover:bg-opacity-90"
                    type="submit"
                >
                    Save
                </button>
            </div>
        </form>
    )
}

export default EditProfileForm