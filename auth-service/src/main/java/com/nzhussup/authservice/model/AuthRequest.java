package com.nzhussup.authservice.model;

import lombok.*;

@NoArgsConstructor
@AllArgsConstructor
@Getter
@Setter
@Builder
public class AuthRequest {

    private String username;
    private String password;
}
