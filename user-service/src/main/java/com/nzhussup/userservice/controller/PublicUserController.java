package com.nzhussup.userservice.controller;

import com.nzhussup.userservice.config.AppConfig;
import com.nzhussup.userservice.dto.PublicUserRegistryRequest;
import com.nzhussup.userservice.dto.UserDTO;
import com.nzhussup.userservice.exceptions.UserAlreadyExistsException;
import com.nzhussup.userservice.exceptions.UserNotFoundException;
import com.nzhussup.userservice.model.User;
import com.nzhussup.userservice.service.UserService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

@RestController
@RequiredArgsConstructor
@RequestMapping(AppConfig.baseApiPath+"public")
public class PublicUserController {

    private final UserService userService;

    @PostMapping
    public ResponseEntity<?> registerUser(@RequestBody PublicUserRegistryRequest publicUser) {
        try {
            UserDTO userDTO = UserDTO.builder()
                    .username(publicUser.getUsername())
                    .password(publicUser.getPassword()).build();

            User savedUser = userService.save(userDTO);
            return ResponseEntity.ok(savedUser);

        } catch (UserAlreadyExistsException e) {
            return ResponseEntity.status(HttpStatus.CONFLICT).body(e.getMessage());
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(e.getMessage());
        }
    }

    @PutMapping
    public ResponseEntity<?> updateUser(@RequestBody User user, Authentication authentication) {
        try {
            UserDTO userDTO = UserDTO.builder()
                    .username(user.getUsername())
                    .password(user.getPassword()).build();
            User updatedUser = userService.update(authentication, user.getId(), userDTO);
            return ResponseEntity.ok(updatedUser);
        } catch (AccessDeniedException e) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).body(e.getMessage());
        } catch (UserNotFoundException e) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(e.getMessage());
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(e.getMessage());
        }
    }

    @DeleteMapping
    public ResponseEntity<?> deleteUser(Authentication authentication, @RequestBody User user) {
        try {
            userService.delete(authentication, user.getId());
            return ResponseEntity.ok().body("User deleted successfully");
        } catch (AccessDeniedException e) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).body(e.getMessage());
        } catch (UserNotFoundException e) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(e.getMessage());
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(e.getMessage());
        }
    }
}
