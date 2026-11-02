package com.nzhussup.baseservice.controller;

import com.nzhussup.baseservice.model.Identifiable;
import com.nzhussup.baseservice.service.BaseService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

public class BaseController<T extends Identifiable> {

    private final BaseService<T> service;

    public BaseController(BaseService<T> service) {
        this.service = service;
    }

    @GetMapping
    public ResponseEntity<List<T>> findAll() {
        List<T> entities = service.findAll();
        return ResponseEntity.ok(entities);
    }

    @PostMapping
    public ResponseEntity<T> save(@RequestBody T entity) {
        try {
            T savedEntity = service.save(entity);
            return ResponseEntity.ok(savedEntity);
        } catch (Exception e) {
            return ResponseEntity.badRequest().body(null);
        }
    }

    @PutMapping("/{id}")
    public ResponseEntity<T> update(@PathVariable Long id, @RequestBody T entity) {
        try {
            if (entity.getId() == null || !id.equals(entity.getId())) {
                return ResponseEntity.badRequest().body(null);
            }
            T updatedEntity = service.update(entity);
            return ResponseEntity.ok(updatedEntity);
        } catch (Exception e) {
            return ResponseEntity.badRequest().body(null);
        }
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(@PathVariable Long id) {
        try {
            service.delete(id);
            return ResponseEntity.ok().build();
        } catch (Exception e) {
            return ResponseEntity.badRequest().body(null);
        }
    }
}
