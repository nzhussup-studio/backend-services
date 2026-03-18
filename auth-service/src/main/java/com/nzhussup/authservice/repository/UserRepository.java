package com.nzhussup.authservice.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import com.nzhussup.authservice.model.User;

public interface UserRepository extends JpaRepository<User, Long> {
    User findByUsername(String username);

}