package com.nzhussup.userservice.repository;

import com.nzhussup.userservice.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface UserRepository extends JpaRepository<User, Long> {

    public User findByUsername(String username);
    public List<User> findByRoleContaining(String role);
}
