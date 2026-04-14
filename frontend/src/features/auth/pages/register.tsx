import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";

import { useMutation } from "@tanstack/react-query";
import { authApi } from "../api/auth.api";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";

const schema = z.object({
    name: z.string().min(2, "Name must be at least 2 characters"),
    email: z.email("Invalid email"),
    password: z.string().min(6, "Minimum 6 characters"),
});

type FormData = z.infer<typeof schema>;

export default function RegisterPage() {
    const navigate = useNavigate();

    const [error, setError] = useState("");

    const mutation = useMutation({
        mutationFn: authApi.register,
        onSuccess: () => navigate("/login"),
        onError: (err: any) => setError(err.message),
    });

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<FormData>({
        resolver: zodResolver(schema),
    });

    const onSubmit = (data: FormData) => {
        try {
            setError("");
            mutation.mutate(data);
            navigate("/login")
        } catch (err: any) {
            setError(err.message);
        }
    };

    return (
        <div className="flex h-screen items-center justify-center px-4">
            <Card className="w-full max-w-md">
                <CardHeader>
                    <CardTitle>Register</CardTitle>
                </CardHeader>

                <CardContent>
                    <form
                        onSubmit={handleSubmit(onSubmit)}
                        className="space-y-4"
                    >
                        <div>
                            <Input
                                placeholder="Name"
                                {...register("name")}
                            />
                            {errors.name && (
                                <p className="text-sm text-red-500">
                                    {errors.name.message}
                                </p>
                            )}
                        </div>

                        <div>
                            <Input
                                placeholder="Email"
                                {...register("email")}
                            />
                            {errors.email && (
                                <p className="text-sm text-red-500">
                                    {errors.email.message}
                                </p>
                            )}
                        </div>

                        <div>
                            <Input
                                type="password"
                                placeholder="Password"
                                {...register("password")}
                            />
                            {errors.password && (
                                <p className="text-sm text-red-500">
                                    {errors.password.message}
                                </p>
                            )}
                        </div>

                        {error && (
                            <p className="text-sm text-red-500">{error}</p>
                        )}

                        <Button
                            type="submit"
                            className="w-full cursor-pointer"
                            disabled={mutation.isPending}
                        >
                            {mutation.isPending
                                ? "Registering..."
                                : "Register"}
                        </Button>
                    </form>

                    <p className="text-sm text-center mt-2">
                        Already have an account?{" "}
                        <Link to="/login" className="text-blue-500">
                            Login
                        </Link>
                    </p>
                </CardContent>
            </Card>
        </div>
    );
}