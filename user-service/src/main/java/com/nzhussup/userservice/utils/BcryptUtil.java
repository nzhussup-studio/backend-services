package com.nzhussup.userservice.utils;

import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;

public class BcryptUtil {

    private static final int STRENGTH = 12;
    private static final BCryptPasswordEncoder passwordEncoder = new BCryptPasswordEncoder(STRENGTH);

    /**
     * Hashes a plaintext password using bcrypt.
     *
     * @param plaintextPassword The plaintext password to hash.
     * @return The hashed password.
     */
    public static String encryptPassword(String plaintextPassword) {
        return passwordEncoder.encode(plaintextPassword);
    }

    /**
     * Validates a plaintext password against a hashed password.
     *
     * @param plaintextPassword The plaintext password.
     * @param hashedPassword    The hashed password.
     * @return True if the password matches, false otherwise.
     */
    public static boolean isPasswordMatch(String plaintextPassword, String hashedPassword) {
        return passwordEncoder.matches(plaintextPassword, hashedPassword);
    }
}
