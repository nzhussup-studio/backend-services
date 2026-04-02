package com.nzhussup.baseservice.service;

import com.nzhussup.baseservice.model.CvGeneratorPreference;
import com.nzhussup.baseservice.repository.CvGeneratorPreferenceRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.LinkedHashMap;
import java.util.List;

@Service
public class CvGeneratorPreferenceService extends BaseService<CvGeneratorPreference> {

    private final CvGeneratorPreferenceRepository preferenceRepository;

    @Autowired
    public CvGeneratorPreferenceService(CvGeneratorPreferenceRepository preferenceRepository) {
        super(preferenceRepository);
        this.preferenceRepository = preferenceRepository;
    }

    @Override
    public List<CvGeneratorPreference> findAll() {
        return findSingleton()
                .map(List::of)
                .orElse(List.of());
    }

    @Override
    public CvGeneratorPreference save(CvGeneratorPreference entity) {
        CvGeneratorPreference preference = entity == null ? new CvGeneratorPreference() : entity;
        if (preference.getPreferencesJson() == null) {
            preference.setPreferencesJson(new LinkedHashMap<>());
        }

        findSingleton().ifPresent(existing -> preference.setId(existing.getId()));
        CvGeneratorPreference saved = preferenceRepository.save(preference);
        removeDuplicates(saved.getId());

        return saved;
    }

    @Override
    public CvGeneratorPreference update(CvGeneratorPreference entity) {
        return save(entity);
    }

    private java.util.Optional<CvGeneratorPreference> findSingleton() {
        List<CvGeneratorPreference> all = preferenceRepository.findAll();
        if (all.isEmpty()) {
            return java.util.Optional.empty();
        }

        CvGeneratorPreference primary = all.getFirst();
        removeDuplicates(primary.getId());
        return java.util.Optional.of(primary);
    }

    private void removeDuplicates(Long keepId) {
        List<CvGeneratorPreference> all = preferenceRepository.findAll();
        for (CvGeneratorPreference item : all) {
            if (keepId == null || !keepId.equals(item.getId())) {
                preferenceRepository.deleteById(item.getId());
            }
        }
    }
}
